# ---------------------------------------------------------------------------
# BUILDER
# ---------------------------------------------------------------------------
FROM --platform=linux/amd64 golang:1.24.4-alpine3.22 AS builder

WORKDIR /build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download && \
    apk --update --no-cache add ca-certificates git

# Copy source code
COPY . .

# Build the Go binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -ldflags="-s -w -X app/api.CommitCount=$(git rev-list --count HEAD) -X app/api.CommitDescribe=$(git describe --always)" \
    -o main ./

# Make entrypoint.sh executable here
COPY docker/builds/lab/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# ---------------------------------------------------------------------------
# FINAL APP
# ---------------------------------------------------------------------------
FROM --platform=linux/amd64 alpine:3.22

WORKDIR /

# Install ca-certificates for Go TLS
RUN apk add --no-cache ca-certificates

VOLUME /public
VOLUME /private

COPY --from=builder /build/main /main
COPY --from=builder /entrypoint.sh /entrypoint.sh
COPY --from=builder /build/public /default_public
COPY --from=builder /build/private /private
COPY --from=builder /build/database/migrations/ /database/migrations
COPY --from=builder /build/database/seeders/ /seeders
COPY --from=builder /build/json_google /json_google
COPY --from=builder /build/docker/builds/lab/.env /.env
COPY --from=builder /build/active.ar.toml /active.ar.toml
COPY --from=builder /build/active.en.toml /active.en.toml

RUN chmod +x /entrypoint.sh

EXPOSE 8000

ENTRYPOINT ["/entrypoint.sh", "/main"]
