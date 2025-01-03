generate:
	go generate -x -run="mockgen" ./...

test:
	go test -v -cover github.com/achu-1612/glcm/hook github.com/achu-1612/glcm/runner github.com/achu-1612/glcm/service