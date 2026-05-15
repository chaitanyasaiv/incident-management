run: 
		go run ./cmd/server/.

build:
		go build -o ./bin/incident-management ./cmd/server/.

.PHONY: run