run:
	go run cmd/main.go
build:
	rm ./argocd-tui
	go build cmd/main.go -o ./argocd-tui
test:
	go test ./...
