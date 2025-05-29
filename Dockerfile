FROM golang:1.24-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o movie-crud-api main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/movie-crud-api .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./movie-crud-api"]
