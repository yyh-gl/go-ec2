.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ./cmd/bin/go-ec2 ./main.go
