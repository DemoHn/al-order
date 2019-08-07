FROM golang:1.12-alpine AS builder

WORKDIR /go/src/app
COPY . . 

RUN mkdir -p /bin

RUN go build -o /bin/alorder ./cmd/main.go

CMD ["/bin/alorder"]