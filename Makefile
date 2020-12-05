.PHONY: all
all: start

# Tools
sync-vendor:
	@echo Download go.mod dependencies
	@go mod download && go mod tidy && go mod vendor

.PHONY: tools
tools:
	echo Installing tools from tools.go
	cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

# GraphQL
graphql: tools
	@echo Graphql Generator
	@gqlgen generate

# Proto
proto: tools
	@go get github.com/golang/protobuf/protoc-gen-go
	@protoc -I ./lib/proto pubsub.proto --go_out=plugins=grpc,paths=source_relative:.

# Development
services-start:
	@echo "starting dev services"
	@docker-compose -p notes-app-dev -f docker-compose.dev.yml up -d

services-stop:
	@echo "stopping dev services"
	@docker-compose -p notes-app-dev -f docker-compose.dev.yml down

services-tail:
	@echo "stopping dev services tail logs"
	@docker-compose -p notes-app-dev -f docker-compose.dev.yml logs -f

.PHONY: start
start: services-start
	@echo "run start app"
	@go run ./cmd/app

.PHONY: stop
stop: services-stop
	@echo "run stop dev"

# Test
test-services-start:
	@echo "starting test services"
	@docker-compose -p notes-app-test -f docker-compose.test.yml up -d

test-services-stop:
	@echo "stopping test services"
	@docker-compose -p notes-app-test -f docker-compose.test.yml down

test-services-tail:
	@echo "stopping test services tail logs"
	@docker-compose -p notes-app-test -f docker-compose.test.yml logs -f

.PHONY: test
test: tools test-services-start
	@echo "run test"
	@ginkgo test -r internal/*

.PHONY: stop-test
stop-test: test-services-stop
	@echo "run stop test"

# Production
build:
	@echo "Prod: build the optimized version for production"
	@CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main ./cmd/app

serve-prod:
	@echo "Prod: serve the app using docker-compose"
	docker-compose up --build

sync-helmfile:
	@echo "Prod: sync helmfile to kubernetes"
	helmfile -f ./deploy/helmfile.yaml sync
