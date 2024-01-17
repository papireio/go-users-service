V := @

NAME := goservice
OUT_DIR := ./target

MAIN_OUT := $(OUT_DIR)/$(NAME)
MAIN_PKG := ./cmd/$(NAME)

default: build

.PHONY: generate
generate:
	 protoc --go_out=./pkg --go_opt=paths=source_relative \
        --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative \
        api/grpc/$(NAME).proto

.PHONY: vendor
vendor:
	$(V)go mod tidy
	$(V)go mod vendor

.PHONY: build
build:
	@echo BUILDING $(MAIN_OUT)
	$(V)go build -o $(MAIN_OUT) $(MAIN_PKG)
	@echo DONE

.PHONY: run
run:
	go run ./cmd/$(NAME)/main.go
