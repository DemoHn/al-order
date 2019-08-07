FROM golang:1.12-alpine AS builder

WORKDIR /go/src/app
COPY . . 

RUN go build -o alorder ./cmd/main.go

FROM alpine:3.7 AS container

COPY --from=builder /go/src/app/alorder /bin/alorder
CMD ["/bin/alorder"]