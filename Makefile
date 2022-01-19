test:
	go test -race -v ./...

cover:
	go test -cover -coverprofile .cover.out .
	go tool cover -html=.cover.out -o coverage.html
	open coverage.html

lint:
	golangci-lint run

.PHONY: test cover
