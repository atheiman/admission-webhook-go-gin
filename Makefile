build:
	go fmt ./api
	dep ensure -v
	env GOOS=linux go build -v -ldflags="-s -w" -o bin/api api/main.go

test:
	go test ./api

run:
	go run api/main.go
