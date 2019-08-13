#!/bin/sh

docker-compose up -d mysql redis
go test ./...