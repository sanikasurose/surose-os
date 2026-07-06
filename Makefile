.PHONY: run build tidy deploy logs ssh stats

DEPLOY_HOST := root@157.245.255.166
DEPLOY_PORT := 2200
DEPLOY_KEY  := ~/.ssh/surose_os_deploy

run:
	go run ./cmd/surose-os

build:
	GOOS=linux GOARCH=amd64 go build -o surose-os ./cmd/surose-os

tidy:
	go mod tidy

deploy:
	GOOS=linux GOARCH=amd64 go build -o surose-os-linux ./cmd/surose-os
	scp -i $(DEPLOY_KEY) -P $(DEPLOY_PORT) surose-os-linux $(DEPLOY_HOST):/home/surose/surose-os/surose-os.new
	ssh -i $(DEPLOY_KEY) -p $(DEPLOY_PORT) $(DEPLOY_HOST) 'systemctl stop surose-os && mv /home/surose/surose-os/surose-os.new /home/surose/surose-os/surose-os && chown surose:surose /home/surose/surose-os/surose-os && chmod +x /home/surose/surose-os/surose-os && systemctl start surose-os'
	rm -f surose-os-linux

logs:
	ssh -i $(DEPLOY_KEY) -p $(DEPLOY_PORT) $(DEPLOY_HOST) journalctl -u surose-os -f

ssh:
	ssh -i $(DEPLOY_KEY) -p $(DEPLOY_PORT) $(DEPLOY_HOST)

stats:
	ssh -i $(DEPLOY_KEY) -p $(DEPLOY_PORT) $(DEPLOY_HOST) "sqlite3 /home/surose/surose-os/data/sessions.db 'SELECT COUNT(*) AS total_visits, COUNT(DISTINCT visitor_hash) AS unique_visitors, CAST(AVG(disconnected_at - connected_at) AS INTEGER) AS avg_session_seconds FROM visits WHERE disconnected_at IS NOT NULL;'"
