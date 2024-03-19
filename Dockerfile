FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build \
    -ldflags "-s -w" \
    -trimpath \
    -o /build/joinfive-bot \
    ./cmd/joinfive-bot/main.go

FROM alpine:3.19
VOLUME ["/data"]
COPY --from=builder /build/joinfive-bot /
EXPOSE 8080
CMD ["/joinfive-bot"]
