.PHONY: start build docs genconfig clean

# Define ENV
BUILD_NAME  := app
BUILD_VER   := $(shell git rev-parse --abbrev-ref HEAD | awk '{gsub(/heads\//,"")}{print $1}' | awk '{gsub(/(\.|\/)/,"-")}{print $1}')
BUILD_TIME  := $(shell date +%s)
BUILD_ARCH  := linux/amd64
BUILD_HASH  := $(shell git rev-parse --short HEAD)
MAIN_PATH   := ./cmd/app/
BUILD_TAG   := ${BUILD_VER}
TARGET_PATH := ./build
BUILD_XPATH := github.com/sendya/pkg/env


all: build

start:
	@go run -tags=doc cmd/app/main.go serve --env dev

build:
	@echo 'Now building ${BUILD_TAG}'
	@env CGO_ENABLED=0 gox -osarch=${BUILD_ARCH} \
		-tags=jsoniter \
		-ldflags="-w -extldflags=-static \
		 -X ${BUILD_XPATH}.Version=${BUILD_VER} \
		 -X ${BUILD_XPATH}.Githash=${BUILD_HASH} \
		 -X ${BUILD_XPATH}.OSArch=${BUILD_ARCH} \
		 -X ${BUILD_XPATH}.Built=${BUILD_TIME} \
		 " \
		-output=${TARGET_PATH}/{{.Dir}}_{{.OS}}_{{.Arch}} ./cmd/app/

docs:
	@swag fmt -d ./cmd/app
	@swag init --pd -g cmd/app/main.go -o third_party/swagger

genconfig:
	@go run cmd/app/main.go genconfig --env dev

clean:
	@rm -rf ./build/
	@go clean -i .