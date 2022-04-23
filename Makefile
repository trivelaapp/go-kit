test: |
	go test ./env/... -cover
	go test ./errors/... -cover
	go test ./http/... -cover
	go test ./log/... -cover
	go test ./metric/... -cover
	go test ./trace/... -cover
