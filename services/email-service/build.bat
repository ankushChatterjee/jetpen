rm -r bin
go build -o bin/email-service.exe cmd/main.go
cp .env bin/