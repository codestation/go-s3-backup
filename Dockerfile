FROM golang:1.21-alpine as builder

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
  ./cmd/go-s3-backup

FROM gitea/gitea:1.21.3 AS gitea
FROM postgres:10-alpine AS postgres-10
FROM postgres:11-alpine AS postgres-11
FROM postgres:12-alpine AS postgres-12
FROM postgres:13-alpine AS postgres-13
FROM postgres:14-alpine AS postgres-14
FROM postgres:15-alpine AS postgres-15
FROM postgres:16-alpine AS postgres-16

FROM alpine:3.19
LABEL maintainer="codestation <codestation@megpoid.dev>"

ENV GITEA_CUSTOM /data/gitea
RUN apk add --no-cache ca-certificates tzdata mariadb-client linux-pam git libpq libedit

COPY --from=gitea /app/gitea /app/gitea
COPY --from=postgres-10 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql10/
COPY --from=postgres-11 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql11/
COPY --from=postgres-12 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql12/
COPY --from=postgres-13 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql13/
COPY --from=postgres-14 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql14/
COPY --from=postgres-15 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql15/
COPY --from=postgres-16 /usr/local/bin/pg_dump /usr/local/bin/pg_restore /usr/local/bin/pg_dumpall /usr/local/bin/psql /usr/libexec/postgresql16/
COPY --from=builder /src/release/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
