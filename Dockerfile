FROM golang:1.23.0-alpine3.20 as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o pingops main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/pingops .
RUN  chmod +x /app/pingops
EXPOSE 8080
CMD ["./pingops"]
