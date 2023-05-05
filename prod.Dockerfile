# Build environment
# -----------------
FROM golang:1.20-alpine as build-env
WORKDIR /rebottle

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api

# Deployment environment
# ----------------------
FROM alpine

COPY --from=build-env /rebottle/bin/api /rebottle/

CMD ["/rebottle/api"]

EXPOSE 8080