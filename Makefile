server:
	go build -o bin/server main.go

clean:
	rm -f bin/*
	
.PHONY: server clean