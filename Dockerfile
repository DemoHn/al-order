FROM golang:1.12-alpine AS builder

RUN apk add --no-cache git

WORKDIR /srv
COPY . .

RUN go build -o alorder ./cmd/main.go

FROM alpine:3.7 AS container

RUN mkdir -p /bin/sql
COPY --from=builder /srv/sql/. /bin/sql
COPY --from=builder /srv/alorder /bin/alorder
CMD ["/bin/alorder"]