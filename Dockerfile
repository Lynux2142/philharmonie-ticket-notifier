FROM golang:1.25.5-alpine
ENV TZ="Europe/Paris"
WORKDIR /app
COPY ./go.mod ./go.sum ./main.go ./web_driver.go ./smtp_client.go ./
RUN go build -o main
CMD ["./main"]
