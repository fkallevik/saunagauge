.PHONY: dev
dev:
	go run main.go -prod=false

.PHONY: run
run:
	go run main.go
