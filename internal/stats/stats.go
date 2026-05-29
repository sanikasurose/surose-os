// Package stats fetches and caches GitHub stats (streak + stars + the full
// contribution calendar) for the home screen, and maintains a persistent
// visitor counter.
//
// Design constraints:
//   - Stats refresh on a background ticker, NOT per keystroke or per request.
//   - Each SSH session takes a Snapshot once at connect time. The home screen
//     reads from its session-local snapshot and never re-fetches.
//   - The visitor counter increments exactly once per session.
//   - All outbound failures degrade gracefully — Snapshot.Available reports
//     whether the GitHub fields are real.
package stats

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

// GridRows × GridCols defines the contribution grid shape: 7 weekdays × 52 weeks.
const (
	GridRows = 7
	GridCols = 52
)

// Snapshot is the per-session view of stats handed to a Bubbletea model.
// Immutable once captured; the model holds it for the life of the session.
type Snapshot struct {
	Visitor            int                       // this session's visitor number (1-indexed)
	Streak             int                       // consecutive days with at least one contribution
	Stars              int                       // total stargazers across owned (non-fork) repos
	TotalContributions int                       // GitHub totalContributions for the last year
	Grid               [GridRows][GridCols]int   // density level 0-4 per cell; [0] = Sunday
	GridStart          time.Time                 // UTC Sunday of the leftmost week column
	Available          bool                      // false if the most recent GitHub fetch failed
	UpdatedAt          time.Time                 // when the cached GitHub data was last refreshed
}

// Cache holds the most recent GitHub data. Safe for concurrent use.
// Run() spawns a background refresher; sessions only read.
type Cache struct {
	mu     sync.RWMutex
	user   string
	token  string
	client *http.Client

	streak             int
	stars              int
	grid               [GridRows][GridCols]int
	gridStart          time.Time
	totalContributions int
	updatedAt          time.Time
	available          bool
}

// NewCache builds a Cache for the given GitHub user. `token` may be empty —
// stars will still be fetched (unauthenticated REST, rate-limited to 60/hr)
// but the contribution calendar (streak + grid + total) requires a token.
func NewCache(user, token string) *Cache {
	return &Cache{
		user:   user,
		token:  token,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// RefreshStatus reports the outcome of a cache refresh. Callers can log
// individual failures without exposing token values.
type RefreshStatus struct {
	StarsErr    error
	CalendarErr error
}

// OK is true when both GitHub fetches succeeded.
func (s RefreshStatus) OK() bool {
	return s.StarsErr == nil && s.CalendarErr == nil
}

// Run does a blocking initial refresh, then starts a background ticker.
// Call once at server startup before accepting SSH sessions — do NOT wrap
// this call in a goroutine; the first refresh must finish first.
func (c *Cache) Run(ctx context.Context, interval time.Duration) RefreshStatus {
	status := c.refresh(ctx)
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				c.refresh(ctx)
			}
		}
	}()
	return status
}

// Snapshot returns the current cached values, stamped with the given visitor
// number. Visitor numbering is the caller's responsibility — pass the result
// of Visitors.Increment().
func (c *Cache) Snapshot(visitor int) Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		Visitor:            visitor,
		Streak:             c.streak,
		Stars:              c.stars,
		TotalContributions: c.totalContributions,
		Grid:               c.grid,
		GridStart:          c.gridStart,
		Available:          c.available,
		UpdatedAt:          c.updatedAt,
	}
}

func (c *Cache) refresh(ctx context.Context) RefreshStatus {
	stars, errS := c.fetchStars(ctx)
	cal, errC := c.fetchCalendar(ctx)

	grid, total, streak, gridStart := computeGrid(cal.days, cal.total)

	c.mu.Lock()
	defer c.mu.Unlock()
	if errS == nil {
		c.stars = stars
	}
	if errC == nil {
		c.streak = streak
		c.grid = grid
		c.gridStart = gridStart
		c.totalContributions = total
	}
	c.updatedAt = time.Now()
	c.available = (errS == nil && errC == nil)
	return RefreshStatus{StarsErr: errS, CalendarErr: errC}
}

func (c *Cache) fetchStars(ctx context.Context) (int, error) {
	total := 0
	for page := 1; page <= 10; page++ {
		url := fmt.Sprintf(
			"https://api.github.com/users/%s/repos?per_page=100&type=owner&page=%d",
			c.user, page,
		)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return 0, err
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}
		resp, err := c.client.Do(req)
		if err != nil {
			return 0, err
		}
		var repos []struct {
			Stargazers int  `json:"stargazers_count"`
			Fork       bool `json:"fork"`
		}
		err = json.NewDecoder(resp.Body).Decode(&repos)
		resp.Body.Close()
		if err != nil {
			return 0, err
		}
		if resp.StatusCode != 200 {
			return 0, fmt.Errorf("stars: GitHub returned %d", resp.StatusCode)
		}
		for _, r := range repos {
			if r.Fork {
				continue
			}
			total += r.Stargazers
		}
		if len(repos) < 100 {
			break
		}
	}
	return total, nil
}

// dayCount is the internal representation of one calendar day.
type dayCount struct {
	Date  time.Time
	Count int
}

// calendarData is the contribution calendar payload from GitHub GraphQL.
type calendarData struct {
	days  []dayCount
	total int // contributionCalendar.totalContributions
}

// fetchCalendar pulls the full year of contribution counts via GraphQL.
// Requires `GITHUB_TOKEN` with `read:user` scope.
//
// Returns days in chronological order (oldest first). The list is typically
// 365–371 days; computeGrid maps them into a 52-week window ending today.
func (c *Cache) fetchCalendar(ctx context.Context) (calendarData, error) {
	if c.token == "" {
		return calendarData{}, fmt.Errorf("calendar: GITHUB_TOKEN not set")
	}

	query := fmt.Sprintf(`{
		user(login: %q) {
			contributionsCollection {
				contributionCalendar {
					totalContributions
					weeks {
						contributionDays {
							date
							contributionCount
						}
					}
				}
			}
		}
	}`, c.user)

	body, _ := json.Marshal(map[string]string{"query": query})
	req, err := http.NewRequestWithContext(
		ctx, "POST", "https://api.github.com/graphql", bytes.NewReader(body),
	)
	if err != nil {
		return calendarData{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return calendarData{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return calendarData{}, fmt.Errorf("calendar: GitHub returned %d", resp.StatusCode)
	}

	var parsed struct {
		Data struct {
			User *struct {
				ContributionsCollection struct {
					ContributionCalendar struct {
						TotalContributions int `json:"totalContributions"`
						Weeks              []struct {
							ContributionDays []struct {
								Date              string `json:"date"`
								ContributionCount int    `json:"contributionCount"`
							} `json:"contributionDays"`
						} `json:"weeks"`
					} `json:"contributionCalendar"`
				} `json:"contributionsCollection"`
			} `json:"user"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return calendarData{}, err
	}
	if len(parsed.Errors) > 0 {
		return calendarData{}, fmt.Errorf("calendar: GraphQL: %s", parsed.Errors[0].Message)
	}
	if parsed.Data.User == nil {
		return calendarData{}, fmt.Errorf("calendar: user %q not found", c.user)
	}

	cal := parsed.Data.User.ContributionsCollection.ContributionCalendar
	var days []dayCount
	for _, w := range cal.Weeks {
		for _, d := range w.ContributionDays {
			t, _ := time.Parse("2006-01-02", d.Date)
			days = append(days, dayCount{Date: t.UTC(), Count: d.ContributionCount})
		}
	}
	if len(days) == 0 {
		return calendarData{}, fmt.Errorf("calendar: no contribution days returned for user %q", c.user)
	}
	return calendarData{days: days, total: cal.TotalContributions}, nil
}

// computeGrid maps calendar days into a 52-week GitHub-style grid ending on the
// current week. Color levels use quartile thresholds of non-zero days, matching
// github.com's relative bucketing. total comes from GitHub's totalContributions.
func computeGrid(days []dayCount, apiTotal int) (grid [GridRows][GridCols]int, total int, streak int, gridStart time.Time) {
	total = apiTotal
	if len(days) == 0 {
		return
	}

	counts := make(map[string]int, len(days))
	var nonZero []int
	for _, d := range days {
		key := d.Date.UTC().Format("2006-01-02")
		counts[key] = d.Count
		if d.Count > 0 {
			nonZero = append(nonZero, d.Count)
		}
	}
	thresholds := contributionThresholds(nonZero)

	today := time.Now().UTC()
	gridStart = sundayOf(today).AddDate(0, 0, -7*(GridCols-1))

	for col := 0; col < GridCols; col++ {
		for row := 0; row < GridRows; row++ {
			day := gridStart.AddDate(0, 0, col*7+row)
			if day.After(today) {
				continue
			}
			grid[row][col] = levelForCount(counts[day.Format("2006-01-02")], thresholds)
		}
	}

	streak = computeStreak(counts, today)
	return
}

// sundayOf returns the UTC Sunday on or before t (GitHub US week start).
func sundayOf(t time.Time) time.Time {
	t = t.UTC()
	return t.AddDate(0, 0, -int(t.Weekday()))
}

// contributionThresholds mirrors GitHub's quartile-based level cutoffs.
func contributionThresholds(nonZero []int) [3]int {
	if len(nonZero) == 0 {
		return [3]int{1, 1, 1}
	}
	sort.Ints(nonZero)
	n := len(nonZero)
	idx := func(p float64) int {
		i := int(float64(n-1) * p)
		if i < 0 {
			return 0
		}
		if i >= n {
			return n - 1
		}
		return i
	}
	return [3]int{
		nonZero[idx(0.25)],
		nonZero[idx(0.50)],
		nonZero[idx(0.75)],
	}
}

func levelForCount(count int, thresholds [3]int) int {
	if count == 0 {
		return 0
	}
	if count <= thresholds[0] {
		return 1
	}
	if count <= thresholds[1] {
		return 2
	}
	if count <= thresholds[2] {
		return 3
	}
	return 4
}

// computeStreak counts consecutive non-zero days ending today (today-with-zero
// does not break the streak, matching GitHub).
func computeStreak(counts map[string]int, today time.Time) int {
	streak := 0
	allowZeroToday := true
	for i := 0; i < 400; i++ {
		d := today.AddDate(0, 0, -i)
		count := counts[d.Format("2006-01-02")]
		if count == 0 {
			if allowZeroToday && i == 0 {
				allowZeroToday = false
				continue
			}
			break
		}
		streak++
		allowZeroToday = false
	}
	return streak
}
