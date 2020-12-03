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
graphql:
	@echo Graphql Generator
	@gqlgen generate

# Development
DOCKER_DEV=postgres-dev-notes-app
DOCKER_DEV_PORT=5432
db-start-dev:
	@echo "starting the postgres docker dev"
	@docker container inspect $(DOCKER_DEV) >/dev/null 2>&1 || docker run --rm --name $(DOCKER_DEV) -e POSTGRES_PASSWORD=postgres -d -p $(DOCKER_DEV_PORT):5432 --mount source=$(DOCKER_DEV),target=/var/lib/postgresql/data postgres:12.3 && ./scripts/wait-for-it.sh localhost:$(DOCKER_DEV_PORT) -- echo "docker dev is up"

db-stop-dev:
	@echo "stopping the postgres docker dev"
	@docker container stop $(DOCKER_DEV)

db-ssh-dev:
	@echo "ssh login the $(DOCKER_DEV) container"
	@docker exec -it $(DOCKER_DEV) psql -U postgres

.PHONY: start
start: db-start-dev
	@echo "run start dev"
	@go run ./cmd/app

# Test
DOCKER_TEST=postgres-test-notes-app
DOCKER_TEST_PORT=5433
db-start-test:
	@echo "starting the postgres docker test"
	@docker container inspect $(DOCKER_TEST) >/dev/null 2>&1 || docker run --rm --name $(DOCKER_TEST) -e POSTGRES_PASSWORD=postgres -d -p 5433:5432 postgres:12.3  && ./scripts/wait-for-it.sh localhost:$(DOCKER_TEST_PORT) -- echo "docker test is up"

db-stop-test:
	@echo "stopping the postgres docker test"
	@docker container stop $(DOCKER_TEST)

db-ssh-test:
	@echo "ssh login the $(DOCKER_TEST) container"
	@docker exec -it $(DOCKER_TEST) psql -U postgres

test: tools db-start-test
	@echo "run test"
	@ginkgo test -r internal/*

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
