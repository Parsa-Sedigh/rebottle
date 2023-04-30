NAME=rebottle

up:
	@echo "Starting the app..."
	cd ./cmd/api && go run .
	@echo "app started!"

build:
	@echo "Building the app..."
	cd ./cmd/api && go build -o ./build/$(NAME)

test:
	@go test -v ./tests/*