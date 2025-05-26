FROM golang:1.22.2 AS builder
WORKDIR /cyderes_app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ingestlog main.go
FROM alpine:latest
WORKDIR /cyderes_app
COPY --from=builder /cyderes_app/ingestlog .
ENTRYPOINT ["/cyderes_app/ingestlog"]
