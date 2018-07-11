FROM golang:1.10-alpine as builder

ARG BUILD_NUMBER=0
ARG COMMIT_SHA
ARG SOURCE_COMMIT
ENV BUILD_COMMIT=${COMMIT_SHA:-${SOURCE_COMMIT:-unknown}}

RUN apk add --no-cache ca-certificates

COPY . $GOPATH/src/megpoid.xyz/go/go-s3-backup/
WORKDIR $GOPATH/src/megpoid.xyz/go/go-s3-backup/

RUN CGO_ENABLED=0 go install -ldflags \
  "-w -s -X main.AppVersion=0.1.${BUILD_NUMBER:-0} -X main.BuildCommit=$(expr substr ${BUILD_COMMIT} 1 8) -X \"main.BuildTime=$(date -u '+%Y-%m-%d %I:%M:%S %Z')\"" \
  -a -tags netgo ./cmd/go-s3-backup

FROM scratch
LABEL maintainer="codestation <codestation404@gmail.com>"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/go-s3-backup /go-s3-backup
COPY --from=builder /tmp /tmp

ENTRYPOINT ["/go-s3-backup"]
