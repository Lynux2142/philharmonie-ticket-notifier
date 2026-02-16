FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY ./go.mod ./go.sum .
RUN go mod tidy
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ="Europe/Paris"
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
