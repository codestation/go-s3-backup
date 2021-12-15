FROM golang:1.17-alpine as builder

ARG CI_COMMIT_TAG
ARG CI_COMMIT_BRANCH
ARG CI_COMMIT_SHA
ARG CI_PIPELINE_CREATED_AT
ARG GOPROXY
ENV GOPROXY=${GOPROXY}

WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/

RUN CGO_ENABLED=0 go build -o release/go-s3-backup \
   -ldflags "-w -s \
   -X main.Version=${CI_COMMIT_TAG:-$CI_COMMIT_BRANCH} \
   -X main.Commit=${CI_COMMIT_SHA:0:8} \
   -X main.BuildTime=${CI_PIPELINE_CREATED_AT}" \
  ./cmd/go-s3-backup

FROM consul:1.10.5 AS consul
FROM gitea/gitea:1.15.7 AS gitea
FROM alpine:3.14
LABEL maintainer="codestation <codestation404@gmail.com>"

ENV GITEA_CUSTOM /data/gitea
RUN apk add --no-cache ca-certificates tzdata postgresql-client mariadb-client linux-pam git

COPY --from=consul /bin/consul /bin/consul
COPY --from=gitea /app/gitea /app/gitea
COPY --from=builder /src/release/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
