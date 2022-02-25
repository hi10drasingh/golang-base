build:
	go build -o server main.go

run: build
		./server

watch:
	~/go/bin/reflex -s -r '\.go$$' make run

