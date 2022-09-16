FROM golang:1.18.6-alpine3.16
RUN apk add build-base

ENV TZ=America/Montreal

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build ./cmd/server/main.go

ENTRYPOINT ["./main"]