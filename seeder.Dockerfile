FROM golang:1.22.2-alpine3.19

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .env

ENTRYPOINT ["tail", "-f", "/dev/null"]
