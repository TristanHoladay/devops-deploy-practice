#!/bin/sh

podman network create ddp-net
podman run -d --name "mongo_server" -p 27017:27017 --network ddp-net mongo:latest
podman run -d --name ddp-api-c -p 8900:8900 -e MONGO_SERVER=mongo_server --network ddp-net ddp-api:latest