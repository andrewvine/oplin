BINARY_NAME=oplin

build:
	GOARCH=amd64 GOOS=darwin go build -o build/${BINARY_NAME}-darwin cmd/oplin/main.go
	GOARCH=amd64 GOOS=linux go build -o build/${BINARY_NAME}-linux cmd/oplin/main.go
	GOARCH=amd64 GOOS=windows go build -o build/${BINARY_NAME}-windows cmd/oplin/main.go

clean:
	go clean
	rm -rf build

test:
	go test -p 1 ./...

init-test-db:
	go run scripts/init.go

run:
	go run cmd/oplin/main.go

test_coverage:
	go test -p 1 ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet ./...

fmt:
	go fmt ./...
