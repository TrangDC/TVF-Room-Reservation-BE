OSFLAG :=
GOARCH :=
VERSION?="1.0.0"
COMMIT?=$(shell git rev-parse --short HEAD)
DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
POSTGRES_URL?="$(POSTGRES_CONNECTION_STRING)"

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OSFLAG = "linux"
	GOARCH = "amd64"
endif
ifeq ($(UNAME_S),Darwin)
	OSFLAG = "darwin"
	GOARCH = "arm64"
endif

ifeq ($(POSTGRES_URL),"")
	POSTGRES_URL="postgres://backend_user:backend_password@localhost:5432/backend_db?sslmode=disable"
endif

os:
	@echo ${OSFLAG}

gengql:
	@go run github.com/99designs/gqlgen generate

genent:
	@go generate ./ent

genmock:
	mockery --dir repositories --all --with-expecter --recursive --inpackage
	mockery --dir internal/azuread --all --with-expecter --recursive --inpackage

gen: gengql genent

build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=$(OSFLAG) GOARCH=$(GOARCH) go build -ldflags "-X main.VERSION=$(VERSION) -X main.COMMIT=$(COMMIT) -X main.DATE=$(DATE) -w -s" -v -o server cmd/main.go

lint:
	@golangci-lint run --fix

test:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@rm -rf coverage.out

hashdb:
	atlas migrate hash --dir file://migrations

migratelint: hashdb
	atlas migrate lint --dev-url $(POSTGRES_URL) --dir file://migrations --latest 2

migratedbdryrun: hashdb
	atlas migrate apply --url $(POSTGRES_URL) --dir file://migrations --dry-run

migratedb: hashdb migratedbdryrun
	atlas migrate apply --url $(POSTGRES_URL) --dir file://migrations

api: build
	./server api

db:
	docker compose up -d db

run:
	docker compose up -d

teardown:
	docker compose down -v --remove-orphans
	docker compose rm --force --stop -v

rebuild:
	docker-compose down -v
	@if [ ! -z "$$(docker images -q)" ]; then docker rmi -f $$(docker images -q); fi
	docker-compose up --build

update-api:
	docker-compose build backend-api
	docker-compose up -d backend-api

format-graphql:
	bash scripts/format-graphql.sh

logapi:
	docker compose logs -f backend-api

initent-%:
	@go run -mod=mod entgo.io/ent/cmd/ent new $*

mockgen-%:
	@mockgen -source=repository/$*.repository.go -destination=repository/mock/mock_$*_repository.go -package=mock
