FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN go build -o /app/server .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]
