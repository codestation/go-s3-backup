FROM golang:1.11-alpine as builder

ARG BUILD_NUMBER=0
ARG COMMIT_SHA
ARG SOURCE_COMMIT
ENV BUILD_COMMIT=${COMMIT_SHA:-${SOURCE_COMMIT:-unknown}}
ENV CGO_ENABLED=0

WORKDIR /app
COPY . /app

RUN go build -mod vendor -ldflags "-w -s \
  -X main.AppVersion=0.1.${BUILD_NUMBER:-0} \
  -X main.BuildCommit=$(expr substr ${BUILD_COMMIT} 1 8) \
  -X \"main.BuildTime=$(date -u '+%Y-%m-%d %I:%M:%S %Z')\"" \
  -o go-s3-backup ./cmd/go-s3-backup

FROM alpine:3.8
LABEL maintainer="codestation <codestation404@gmail.com>"

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
