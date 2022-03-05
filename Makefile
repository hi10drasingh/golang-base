build:
	go build -o server main.go

run: build
		./server

lint:
	golangci-lint run

watch:
	~/go/bin/reflex -s -r '\.go$$' make run

