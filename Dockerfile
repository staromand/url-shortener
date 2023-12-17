FROM golang:1.21.5

ENV GO111MODULE=on

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o ./cmd/url-shortener/bin ./cmd/url-shortener
ENTRYPOINT ["./cmd/url-shortener/bin"]