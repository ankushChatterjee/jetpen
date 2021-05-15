FROM golang:1.16.4-alpine
WORKDIR /src/app

COPY . ./
RUN go mod download
RUN go build -o bin/newsletter-service cmd/main.go
EXPOSE 3000
CMD ["./bin/newsletter-service"]