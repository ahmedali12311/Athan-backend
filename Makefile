include docker/builds/local/.env
export

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
# ============================================================================= 
# Go
# ============================================================================= 
.PHONY: list update
## list: check installed libraries for the project
list:
	go list -m -u
## update: updates installed libraries for the project
update:
	go get -u ./...

# ============================================================================= 
# Local Dev
# ============================================================================= 
## air: runs the project locally using air docker image
air:
	@set -a && source docker/builds/local/.env
	@docker run -it --rm \
		--network ${NETWORK} \
		--env-file docker/builds/local/.env \
		-w "/${APP_CODE}" \
		-e "air_wd=/${APP_CODE}" \
		-v ${ROOT_DIR}:/${APP_CODE} \
		-v ~/go/pkg/mod:/go/pkg/mod \
		-p ${PORT}:${PORT} \
		cosmtrek/air \
		-c ./.air.toml


# ============================================================================= 
# Translations
# ============================================================================= 
.PHONY: translate/extract translate translate/done
## translate: extracts new messages and prepare for translations
translate: translate/extract
	@goi18n extract --outdir .
	@goi18n merge --outdir . ./active.ar.toml ./active.en.toml
## translate/done: merges new translations into active.*.toml files
translate/done:
	@goi18n merge --sourceLanguage ar --outdir . ./active.ar.toml ./translate.ar.toml


# ============================================================================= 
# Migrations
# ============================================================================= 
/PHONY: migrate/up migrate/up/all migrate/down migrate/down/all migrate/force
## migrate/up n=<number>: migrates up n steps
migrate/up:
	docker run --rm -v ./database/migrations:/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up $(n)
## migrate/up/all: migrates up to latest
migrate/up/all:
	docker run --rm -v ./database/migrations:/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up
## migrate/down n=<number>: migrates down n steps
migrate/down:
	docker run --rm -v ./database/migrations:/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down $(n)
## migrate/down/all: migrates down all steps
migrate/down/all:
	docker run --rm -v ./database/migrations:/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down -all
## migration n=<file_name>: creates migration files up/down for file_name
migration:
	docker run --rm -v ./database/migrations:/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ create -seq -ext=.sql -dir=./migrations $(n)
## migrate/force n=<version>: forces migration version number
migrate/force:
	docker run --rm -v ./database/migrations:/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database=$(CONNECTION_STRING) force $(n)


# ============================================================================= 
# STANDALONE SERVER COMMANDS
# ============================================================================= 
.PHONY: certs renew reload
## certs: checks current certificates by certbot
certs:
	docker compose run --rm  certbot certificates
## renew: requests a renewal for ssl certificate by certbot
renew:
	docker compose run --rm  certbot renew
## reload: reloads nginx configuration inside container on server
reload:
	docker exec -it nginx nginx -s reload
## request: creates a certificate for subdomain provided in (n)
##        : eg > make request n=blueprint.sadeem-lab.com
request:
	docker compose run --rm certbot certonly --webroot --webroot-path /var/www/certbot/ -d $(n)

## check-space: detailed disk usage of / showing only Mega and Giga byte dirs
check-space:
	du -cha --max-depth=1 / | grep -E "M|G"
# Docker
.PHONY: prune ps inspect
## prune: clears system volumes
prune:
	docker system prune -a -f --volumes
## ps: lists docker containers in a formatted style
ps:
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.RunningFor}}\t{{.Size}}\t{{.Ports}}"
## inspect: a container local ip n=name of container
inspect: 
	docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" $(n)
## prune-dangled-volumes: clear system dangled volumes
prune-dangled-volumes:
	docker volume ls -q -f dangling=true | xargs -r docker volume rm
