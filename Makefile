include docker/builds/local/.env
export

prune:
	docker system prune -a -f --volumes
ps:
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.RunningFor}}\t{{.Size}}\t{{.Ports}}"
inspect: 
	docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" $(n)

# Go
.PHONY: list update
list:
	go list -m -u
update:
	go get -u ./...

air:
	docker run -it --rm \
		--network ${NETWORK} \
		-w "/${APP_CODE}" \
		-e "air_wd=/${APP_CODE}" \
		-v ${ROOT_DIR}:/${APP_CODE} \
		-v ~/go/pkg/mod:/go/pkg/mod \
		-p ${PORT}:${PORT} \
		cosmtrek/air
		-c ./.air.toml
