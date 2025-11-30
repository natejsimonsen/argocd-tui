run:
	go run cmd/main.go
build:
	rm -f ./argocd-tui
	go build -o ./argocd-tui cmd/main.go 
test:
	go test ./...
