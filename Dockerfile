FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN  go build -o mail_app ./cmd/mail

FROM alpine:3.23
WORKDIR /root/
RUN apk --no-cache add ca-certificates
EXPOSE 8080
COPY --from=builder /app/mail_app .
ENTRYPOINT ["./mail_app"]