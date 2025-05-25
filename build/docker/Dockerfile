FROM golang:1.23 AS go_image
WORKDIR /cyderes_app
COPY go.mod ./
# RUN go mod download
COPY ./app/backend/. .
#RUN go build -o main .
#FROM alpine:latest
#WORKDIR /root/
#COPY --from=go_image /cyderes_app/main .
CMD ["go", "run", "."]
