all: start
.PHONY: all

tools:
	@echo "install tools"
	command -v ginkgo >/dev/null 2>&1 || go install github.com/onsi/ginkgo/ginkgo

db-start-dev:
	@echo "starting the postgres docker dev"
	docker container inspect docker-postgres-dev >/dev/null 2>&1 || docker run --rm --name docker-postgres-dev -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 -v ${DEVROOT}/docker/volumes/postgres:/var/lib/postgresql/data postgres:12.3
	docker exec -it docker-postgres-dev psql -U postgres -c "CREATE DATABASE notes_app_dev;" >/dev/null 2>&1 || true

db-stop-dev:
	@echo "stopping the postgres docker dev"
	docker container stop docker-postgres-dev

db-ssh-dev:
	@echo "ssh login the docker-postgres container"
	docker exec -it docker-postgres-dev psql -U postgres

start: db-start-dev
	@echo "run start dev"
	go run main.go


db-start-test:
	@echo "starting the postgres docker test"
	docker container inspect docker-postgres-test >/dev/null 2>&1 || docker run --rm --name docker-postgres-test -e POSTGRES_PASSWORD=postgres -d -p 5433:5432 postgres:12.3
	docker exec -it docker-postgres-test psql -U postgres -c "CREATE DATABASE notes_app_test;" >/dev/null 2>&1 || true

db-stop-test:
	@echo "stopping the postgres docker test"
	docker container stop docker-postgres-test

db-ssh-test:
	@echo "ssh login the docker-postgres-test container"
	docker exec -it docker-postgres-test psql -U postgres

test: tools db-start-test
	@echo "run test"
	ginkgo test -r *