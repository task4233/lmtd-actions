run:
	go run cmd/lmtd-actions/main.go

test:
	go test -cover -v -coverprofile=cover.out ./   
	go tool cover -html=cover.out -o ./cover.html
	open cover.html