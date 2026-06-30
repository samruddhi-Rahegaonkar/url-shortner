FROM golang:1.25.6 AS builder
WORKDIR /urlshortner
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /urlshortner/server .
EXPOSE 8080
CMD ["./server"]