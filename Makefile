all: start
.PHONY: all

# Tools
sync-vendor:
	@echo Download go.mod dependencies
	@go mod download && go mod tidy && go mod vendor

tools:
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

# Development
db-start-dev:
	@echo "starting the postgres docker dev"
	@mkdir -p ${DEVROOT}/docker/volumes/postgres
	@docker container inspect docker-postgres-dev >/dev/null 2>&1 || docker run --rm --name docker-postgres-dev -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 -v ${DEVROOT}/docker/volumes/postgres:/var/lib/postgresql/data postgres:12.3

db-stop-dev:
	@echo "stopping the postgres docker dev"
	@docker container stop docker-postgres-dev

db-ssh-dev:
	@echo "ssh login the docker-postgres container"
	@docker exec -it docker-postgres-dev psql -U postgres

start: db-start-dev
	@echo "run start dev"
	@go run main.go

# Test
db-start-test:
	@echo "starting the postgres docker test"
	@docker container inspect docker-postgres-test >/dev/null 2>&1 || docker run --rm --name docker-postgres-test -e POSTGRES_PASSWORD=postgres -d -p 5433:5432 postgres:12.3

db-stop-test:
	@echo "stopping the postgres docker test"
	@docker container stop docker-postgres-test

db-ssh-test:
	@echo "ssh login the docker-postgres-test container"
	@docker exec -it docker-postgres-test psql -U postgres

test: tools db-start-test
	@echo "run test"
	@ginkgo test -r *

# Production
build:
	@echo "Prod: build the optimized version for production"
	@CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

serve-prod:
	@echo "Prod: serve the app using docker-compose"
	#docker volume create notes-app-postgres-data
	docker-compose -f ./deploy/docker-compose.yml up --build
