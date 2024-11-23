FROM golang:1.23-alpine AS builder

WORKDIR /build
COPY go.mod .
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app app/main.go

FROM alpine:latest
WORKDIR /app
EXPOSE 8000
COPY --from=builder app /bin/main
CMD ["/bin/main"]
