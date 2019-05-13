FROM golang:1.12-alpine as builder

ARG BUILD_NUMBER
ARG BUILD_COMMIT_SHORT
ENV GO111MODULE on
ENV CGO_ENABLED 0
WORKDIR /src

COPY . .

RUN go build -o release/go-s3-backup \
   -mod vendor -ldflags "-w -s \
  -X main.AppVersion=0.1.${BUILD_NUMBER} \
  -X main.BuildCommit=${BUILD_COMMIT_SHORT} \
  -X \"main.BuildTime=$(date -u '+%Y-%m-%d %I:%M:%S %Z')\"" \
  ./cmd/go-s3-backup

FROM consul:1.5 AS consul
FROM gogs/gogs:0.11.86 AS gogs
FROM alpine:3.9
LABEL maintainer="codestation <codestation404@gmail.com>"

ENV GOGS_CUSTOM /data/gogs
RUN apk add --no-cache ca-certificates tzdata postgresql-client mariadb-client linux-pam

COPY --from=consul /bin/consul /bin/consul
COPY --from=gogs /app/gogs /app/gogs
COPY --from=builder /src/release/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
