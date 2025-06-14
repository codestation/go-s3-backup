FROM golang:1.24-alpine AS builder

ARG CI_COMMIT_TAG
ARG GOPROXY
ENV GOPROXY=${GOPROXY}

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/

RUN set -ex; \
    CGO_ENABLED=0 go build -o release/go-s3-backup \
    -trimpath \
    -ldflags "-w -s \
   -X version.Tag=${CI_COMMIT_TAG}" \
  .

FROM postgres:12.22-alpine AS postgres-12
FROM postgres:13.21-alpine AS postgres-13
FROM postgres:14.18-alpine AS postgres-14
FROM postgres:15.13-alpine AS postgres-15
FROM postgres:16.9-alpine AS postgres-16
FROM postgres:17.5-alpine AS postgres-17

FROM alpine:3.22
LABEL maintainer="codestation <codestation@megpoid.dev>"

RUN apk add --no-cache ca-certificates tzdata mariadb-client libpq libedit zstd-libs lz4-libs

COPY --from=postgres-12 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql12/
COPY --from=postgres-13 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql13/
COPY --from=postgres-14 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql14/
COPY --from=postgres-15 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql15/
COPY --from=postgres-16 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql16/
COPY --from=postgres-17 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql17/

COPY --from=builder /src/release/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
