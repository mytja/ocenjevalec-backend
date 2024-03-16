FROM golang:1.18-alpine AS builder

COPY . /app

WORKDIR /app

# Add gcc
RUN apk add build-base

RUN go mod download && \
    go env -w GOFLAGS=-mod=mod && \
    go get . && \
    go build -v -o backend .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/backend ./backend

EXPOSE 80
CMD [ "./backend", "--useenv" ]
