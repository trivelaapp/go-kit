#!make

GOPATH=$(shell go env GOPATH)

run_client:
	go run client/example.go | jq

run_server:
	GIN_MODE=release go run server/example.go | jq
