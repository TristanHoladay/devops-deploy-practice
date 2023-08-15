FROM golang:1.19

WORKDIR /build
COPY go.mod go.sum ./
COPY .env ./.env
COPY src/frontend/index.html /build/ui/
COPY src/backend ./src/backend
RUN GO111MODULE=on CGO_ENABLED=0 go build -o /build/ddp-api ./src/backend/main.go
EXPOSE 8900
CMD ["./ddp-api"]