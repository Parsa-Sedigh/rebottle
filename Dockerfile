FROM golang:1.20-alpine
WORKDIR /rebottle

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api

CMD ["/rebottle/bin/api"]

EXPOSE 8080