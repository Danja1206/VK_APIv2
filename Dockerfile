FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/vkProject

EXPOSE 8080

CMD ["./main"]