DEBUG=0

POSTGRES_PASSWORD=crc
POSTGRES_USER=crc
POSTGRES_DB=gumbaroo
POSTGRES_PORT=5432:5432
POSTGRES_HOST=localhost

.PHONY: build

build: platform-changelog-api platform-changelog-migration

platform-changelog-api:
	go build -o $@ cmd/api/main.go

platform-changelog-migration:
	go build -o $@ internal/migration/main.go

lint:

	gofmt -l .
	gofmt -s -w .

test:

	go test -p 1 -v ./...

run-migration: platform-changelog-migration

	./platform-changelog-migration

run-api: platform-changelog-api

	DEBUG=${DEBUG} ./platform-changelog-api

run-api-mock: platform-changelog-api

	DEBUG=${DEBUG} DB_IMPL=mock ./platform-changelog-api

run-db:

	podman run --rm -it -p ${POSTGRES_PORT} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_DB=${POSTGRES_DB} --name postgres postgres:12.4

check-db:

	psql -h ${POSTGRES_HOST} --user ${POSTGRES_USER} --db ${POSTGRES_DB}

test-github-webhook:

	curl -X POST -H "X-Github-Event: push" -H "Content-Type: application/json" --data "@tests/github_webhook.json" http://localhost:8000/api/platform-changelog/v1/github-webhook

test-gitlab-webhook:

	curl -X POST -H "X-Gitlab-Event: Push Hook" -H "Content-Type: application/json" --data "@tests/gitlab_webhook.json" http://localhost:8000/api/platform-changelog/v1/gitlab-webhook

test-tekton-task:

	curl -X POST http://localhost:8000/api/platform-changelog/v1/tekton --data "@tests/tekton/valid.json" -H "Content-Type: application/json"

compose:

	podman-compose -f development/compose.yml up

compose-quiet:

	podman-compose -f development/compose.yml up -d

compose-down:

	podman-compose -f development/compose.yml down

clean:
	go clean
	rm -f platform-changelog-api platform-changelog-migration
