FROM golang:1.14-alpine as builder

ARG CI_TAG
ARG BUILD_NUMBER
ARG BUILD_COMMIT_SHORT
ARG CI_BUILD_CREATED
ENV GO111MODULE on
ENV CGO_ENABLED 0
WORKDIR /src

COPY . .

RUN go build -o release/go-s3-backup \
   -ldflags "-w -s \
   -X main.Version=${CI_TAG} \
   -X main.BuildNumber=${BUILD_NUMBER} \
   -X main.Commit=${BUILD_COMMIT_SHORT} \
   -X main.BuildTime=${CI_BUILD_CREATED}" \
  ./cmd/go-s3-backup

FROM consul:1.8 AS consul
FROM gitea/gitea:1.12 AS gitea
FROM alpine:3.12
LABEL maintainer="codestation <codestation404@gmail.com>"

ENV GITEA_CUSTOM /data/gitea
RUN apk add --no-cache ca-certificates tzdata postgresql-client mariadb-client linux-pam

COPY --from=consul /bin/consul /bin/consul
COPY --from=gitea /app/gitea /app/gitea
COPY --from=builder /src/release/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
