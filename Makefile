all: start
.PHONY: all

tools:
	@echo "install tools"
	command -v ginkgo >/dev/null 2>&1 || go install github.com/onsi/ginkgo/ginkgo

db-start:
	@echo "starting the postgres docker"
	docker container inspect docker-postgres >/dev/null 2>&1 || docker run --rm --name docker-postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 -v ${DEVROOT}/docker/volumes/postgres:/var/lib/postgresql/data postgres:12.3
	docker exec -it docker-postgres psql -U postgres -c "CREATE DATABASE notes_app_test;" >/dev/null 2>&1 || true
	docker exec -it docker-postgres psql -U postgres -c "CREATE DATABASE notes_app_dev;" >/dev/null 2>&1 || true

db-stop:
	@echo "stopping the postgres docker"
	docker container stop docker-postgres

db-ssh:
	@echo "ssh login the docker-postgres container"
	docker exec -it docker-postgres psql -U postgres

test: tools db-start
	@echo "run test"
	ginkgo test -r *

start: db-start
	@echo "run start dev"
	go run main.go

