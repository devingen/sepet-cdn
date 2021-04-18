# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:latest

LABEL maintainer="Emir Luleci <emir@devingen.io>"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o sepet-cdn cmd/sepet-cdn/sepet-cdn.go

EXPOSE 80

CMD ["./sepet-cdn"]