FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY app/go.mod .
RUN go mod download
COPY app/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM alpine:3.20 
WORKDIR /app
RUN adduser -D -g appuser appuser && chown -R appuser:appuser /app
COPY --from=builder /app/main .
USER appuser
EXPOSE 8080
CMD ["./main"]
