APP_NAME = basecoat
GIT_COMMIT := $(shell git rev-parse --short HEAD)
SHELL = /bin/bash
VERSION=$(shell date +%s)


GO_LDFLAGS := '-X "github.com/clintjedwards/${APP_NAME}/cmd.appVersion=$(VERSION) $(GIT_COMMIT)" \
			   -X "github.com/clintjedwards/${APP_NAME}/service.appVersion=$(VERSION) $(GIT_COMMIT)"'

## build: run tests and compile full app in production mode
build: export FRONTEND_API_HOST="https://${APP_NAME}.clintjedwards.com"
build: check-path-included
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
	go mod tidy
	go test ./utils
	npm run --prefix ./frontend build:production
	packr build -ldflags $(GO_LDFLAGS) -o $(path)

## build-backend: build backend without frontend assets
build-backend: check-path-included
	protoc --go_out=plugins=grpc:. api/*.proto
	go mod tidy
	go test ./utils
	go build -ldflags $(GO_LDFLAGS) -o $(path)

## build-dev: build development version of app
build-dev:
	npx webpack --config="./frontend/webpack.config.js" --mode="development"
	packr build -ldflags $(GO_LDFLAGS) -o $(path)

## build-protos: build required protobuf files
build-protos:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto

## help: prints this help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## install: build application and install on system
install:
	protoc --go_out=plugins=grpc:. api/*.proto
	go mod tidy
	npm run --prefix ./frontend build:production
	packr build -ldflags $(GO_LDFLAGS) -o /tmp/${APP_NAME}
	sudo mv /tmp/${APP_NAME} /usr/local/bin/
	chmod +x /usr/local/bin/${APP_NAME}

## run: build application and run server; useful for dev
run: export FRONTEND_API_HOST="https://localhost:8080"
run:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
	go mod tidy
	npm run --prefix ./frontend build:development
	packr build -ldflags $(GO_LDFLAGS) -o /tmp/${APP_NAME} && /tmp/${APP_NAME} server

## run-backend: build backend only and run server; useful for dev
run-backend:
	protoc --go_out=plugins=grpc:. api/*.proto
	go mod tidy
	packr build -ldflags $(GO_LDFLAGS) -o /tmp/${APP_NAME} && /tmp/${APP_NAME} server

check-path-included:
ifndef path
	$(error path is undefined; ex. path=/tmp/${APP_NAME})
endif
