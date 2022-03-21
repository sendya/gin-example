.PHONY: start build genconfig clean

# ENV
BUILD_VER   := $(shell git rev-parse --abbrev-ref HEAD | awk '{gsub(/heads\//,"")}{print $1}' | awk '{gsub(/(\.|\/)/,"-")}{print $1}')
BUILD_TIME  := $(shell date +%s)
BUILD_ARCH  := linux/amd64
BUILD_HASH  := $(shell git rev-parse --short HEAD)
MAIN_PATH   := ./cmd/app/
BUILD_TAG   := ${BUILD_VER}

all: build

start:
	@go run cmd/app/main.go -env dev

build:
	@echo 'Now building ${BUILD_NAME}:${BUILD_TAG}'

genconfig:
	@go run cmd/app/main.go -env "" -genconfig

clean:
	rm -rf ./build/
	go clean -i .