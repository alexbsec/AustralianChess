FROM golang:1.26.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

ENV PATH="/root/go/bin:${PATH}"

EXPOSE 8080

CMD ["go", "run", "./cmd/api"]
