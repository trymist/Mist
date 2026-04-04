FROM golang:1.25-alpine

ENV CGO_ENABLED=1
ENV GO111MODULE=on
ENV PATH="/root/go/bin:${PATH}"
ENV ENV="dev"

RUN apk add --no-cache \
    build-base \
    git \
    curl \
    sqlite \
    bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o myapp ./...

EXPOSE 8080

CMD ["./myapp"]
