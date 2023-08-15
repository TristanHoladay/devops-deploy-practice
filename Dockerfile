FROM golang:1.19 as build2

WORKDIR /build
COPY go.mod go.sum ./
COPY src/backend ./src/backend
RUN GO111MODULE=on CGO_ENABLED=0 go build -o /build/form-api ./src/backend/main.go
