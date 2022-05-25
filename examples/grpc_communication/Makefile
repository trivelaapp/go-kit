#!make

GOPATH=$(shell go env GOPATH)

run_client:
	go run client/main.go | jq

run_server:
	go run server/main.go | jq
