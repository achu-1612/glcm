generate:
	go generate -x -run="mockgen" ./...

test:
	go test -v -cover ./...