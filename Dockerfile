FROM golang:1.22-alpine AS builder
RUN apk add --no-cache --update gcc g++
WORKDIR /src
COPY go.mod go.sum ./
ENV CGO_ENABLED=1
RUN go mod download
COPY . .
RUN go build -o /build/joinfive-bot ./cmd/joinfive-bot/main.go

FROM alpine:3.19
VOLUME ["/data"]
COPY --from=builder /build/joinfive-bot /
EXPOSE 8080
CMD ["/joinfive-bot"]
