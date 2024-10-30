# BUILDER ---------------------------------------------------------------------
FROM --platform=linux/amd64 golang:1.22.2-alpine3.19 AS builder

WORKDIR /build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download && \ 
    apk --update --no-cache add ca-certificates git

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -ldflags="-s -w -X app/api.CommitCount=$(git rev-list --count HEAD) -X app/api.CommitDescribe=$(git describe --always)" \
    -o main ./

# FINAL APP -------------------------------------------------------------------
FROM --platform=linux/amd64 scratch

WORKDIR /

VOLUME /public
VOLUME /private

COPY --from=builder ["/build/main", "/main"]
COPY --from=builder ["/build/public", "/public"]
COPY --from=builder ["/build/private", "/private"]
COPY --from=builder ["/build/database/migrations/", "/migrations"]
COPY --from=builder ["/build/database/seeders/", "/seeders"]
COPY --from=builder ["/build/json_google", "/json_google"]
COPY --from=builder ["/build/docker/builds/local/.env", "/.env"]
COPY --from=builder ["/build/active.ar.toml", "/active.ar.toml"]
COPY --from=builder ["/build/active.en.toml", "/active.en.toml"]
COPY --from=builder ["/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/ca-certificates.crt"]

EXPOSE 8000

ENTRYPOINT ["/main"]
