FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/pvz_server ./cmd/apiserver/main.go

EXPOSE 8080

CMD ["./bin/pvz_server"]