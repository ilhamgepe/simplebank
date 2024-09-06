# STAGE 1: build the go application
FROM golang:1.23.0-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


# STAGE 2: create the final image
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["/app/main"]