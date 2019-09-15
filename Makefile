SHELL := /bin/bash

install:
	go get

start:
	make build
	GO_ENV=development go run *.go

build:
	npm run build

lint:
	$$(npm bin)/eslint . --ignore-path .gitignore --cache