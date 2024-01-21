V := @

NAME := go-users
OUT_DIR := ./target

MAIN_OUT := $(OUT_DIR)/$(NAME)
MAIN_PKG := ./cmd/$(NAME)

default: build

.PHONY: build
build:
	@echo BUILDING $(MAIN_OUT)
	$(V)go build -o $(MAIN_OUT) $(MAIN_PKG)
	@echo DONE

.PHONY: run
run:
	go run ./cmd/$(NAME)/main.go

.PHONY: vendor
vendor:
	$(V)go mod tidy -go=1.21 && go mod tidy -go=1.21
	$(V)go mod tidy
	$(V)go mod vendor

.PHONY: proto
proto:
	 protoc --go_out=./pkg --go_opt=paths=source_relative \
        --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative \
        api/grpc/$(NAME).proto
