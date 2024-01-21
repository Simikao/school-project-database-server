server:
	go build -o bin/server main.go

clean:
	rm -f bin/*
	
dup:
	docker compose up -d

ddown:
	docker compose down

.PHONY: server clean dup ddown