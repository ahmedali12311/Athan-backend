include .env
export

.PHONY: init
 init:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/fraenky8/tables-to-go@master
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
	@go install github.com/jesseduffield/lazygit@latest
	@go install github.com/mfridman/tparse@latest
	@go install github.com/nametake/golangci-lint-langserver@latest
	@go install github.com/nicksnyder/go-i18n/v2/goi18n@latest
	@go install github.com/segmentio/golines@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install mvdan.cc/gofumpt@latest

# Application
.PHONY: run test test-race
run:
	@air
test:
	@go test main_test.go -json | tparse -all
test-race:
	@go test -cover -race main_test.go -json | tparse -all

# Translations
.PHONY: translate.extract translate.merge
translate.extract:
	@goi18n extract --outdir .
translate.merge:
	@goi18n merge --outdir . ./active.ar.toml ./active.en.toml
translate.merge.done:
	@goi18n merge --sourceLanguage ar --outdir . ./active.ar.toml ./translate.ar.toml

# Migrations
.PHONY: migrate.up migrate.up.all migrate.down migrate.down.all migrate.force refresh
migrate.up:
	docker run --rm -v ./$(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up $(n)
migrate.up.all:
	docker run --rm -v ./$(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up
migrate.down:
	docker run --rm -v ./$(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down $(n)
migrate.down.all:
	docker run --rm -v ./$(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down -all
migration:
	docker run --rm -v ./$(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ create -seq -ext=.sql -dir=./migrations $(n)
migrate.force:
	docker run --rm -v ./$(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database=$(CONNECTION_STRING) force $(n)
refresh: migrate.down.all migrate.up seed

# Seeders
.PHONY: seed, live-seed, live-seed.down
seed:
	@go run $(SEEDERS_ROOT)/...	
live-seed:
	docker compose -f docker-compose-seeder.yml up --build -d
live-seed.down:
	docker compose -f docker-compose-seeder.yml down

CONTAINER_NAME:=${CONTAINER_NAME}
CONTAINER_TAG:=${APP_VER}.$(shell git rev-list --count HEAD).$(shell git describe --always)
CONTAINER_IMG:=${CONTAINER_NAME}:${CONTAINER_TAG}

# hub.docker.com: Production
.PHONY: dh dh/down dh/local dh/local/down dh/push db/conn volumes prune ps inspect
dh:
	export CONTAINER_TAG=${CONTAINER_TAG}
	docker compose -f docker-compose-dh.yml up --build -d
dh/down:
	docker compose -f docker-compose-dh.yml down

dh/local:
	export CONTAINER_TAG=${CONTAINER_TAG}
	docker compose -f docker-compose-dh-local.yml up --build -d
dh/local/down:
	docker compose -f docker-compose-dh-local.yml down

dh/push: dh
	docker tag sadeem/${CONTAINER_IMG} ${CONTAINER_REG}/${CONTAINER_IMG}
	docker tag sadeem/${CONTAINER_IMG} ${CONTAINER_REG}/${CONTAINER_NAME}:latest
	docker push ${CONTAINER_REG}/${CONTAINER_NAME} -a
db/conn:
	psql ${CONNECTION_STRING}
volumes:
	docker volume create pale_skull_public && \
    docker volume create pale_skull_private
prune:
	docker system prune -a -f --volumes
ps:
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.RunningFor}}\t{{.Size}}\t{{.Ports}}"
# inspect a container local ip n=name of container
inspect: 
	docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" $(n)

# Go
.PHONY: list update
list:
	go list -m -u
update:
	go get -u ./...

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	gofumpt -l -w -extra .
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# psql/
psql:
	docker run -it --rm --network ${NETWORK} postgis/postgis psql -h ${PG_HOST} -U ${PG_USER}

# build locally 
build:
	go build -o main .
# build with tags
build/local:
	go build -tags local -o main . 