SVC_NAME := bidder
APP_EXECUTABLE := ./bidder-svc

.PHONY: all
all: build

.PHONY: build
build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -o $(APP_EXECUTABLE) --tags netgo --ldflags '-extldflags "-lm -lstdc++ -static"' ./cmd/bidder
	docker rmi -f $(SVC_NAME)
	docker build --build-arg APP_PORT -t $(SVC_NAME) .