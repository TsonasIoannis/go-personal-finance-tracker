go mod tidy
GOOS=linux GOARCH=amd64 go build -o dist/personal-finance-tracker-linux cmd/main.go
GOOS=darwin GOARCH=amd64 go build -o dist/personal-finance-tracker-mac cmd/main.go
GOOS=windows GOARCH=amd64 go build -o dist/personal-finance-tracker.exe cmd/main.go