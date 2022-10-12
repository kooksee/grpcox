build:
	@go build

run:
	@go run .

usage: FORCE
	exit 1

FORCE:

include config.env
export $(shell sed 's/=.*//' config.env)

start: FORCE
	@echo " >> building..."
	@mkdir -p log
	@go build
	@./grpcox

.PHONY: app
app:
	GOARCH=wasm GOOS=js go build -v -o web/app.wasm ./web/main/*
	go run .