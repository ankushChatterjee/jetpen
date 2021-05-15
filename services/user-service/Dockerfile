FROM golang:1.16.4-alpine
WORKDIR /src/app

COPY . ./
RUN go mod download
RUN go build -o bin/user-service cmd/main.go
CMD ["./bin/user-service"]