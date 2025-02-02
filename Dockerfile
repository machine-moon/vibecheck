FROM golang:1.23.5-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 8080
COPY secrets.yml .
RUN source secrets.yml
CMD ["./server"]

