FROM golang:1.16.4-alpine
WORKDIR /src/app

COPY . ./
RUN go mod download
RUN go build -o bin/sub-service cmd/main.go
CMD ["./bin/sub-service"]