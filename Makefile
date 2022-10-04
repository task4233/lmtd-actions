build:
	go build -o lmtd-checker ./cmd/lmtd-checker/main.go

run:
	go run cmd/lmtd-checker/main.go

test:
	go test -race -shuffle=on -v ./...

test-with-coverage:
	go test -race -shuffle=on -cover -v -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o ./cover.html
	open cover.html