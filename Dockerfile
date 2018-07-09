FROM golang:1.10-alpine as builder

ARG BUILD_NUMBER=0
ARG COMMIT_SHA
ARG SOURCE_COMMIT
ENV BUILD_COMMIT=${COMMIT_SHA:-${SOURCE_COMMIT:-unknown}}

COPY . $GOPATH/src/megpoid.xyz/go/go-s3-backup/
WORKDIR $GOPATH/src/megpoid.xyz/go/go-s3-backup/

RUN CGO_ENABLED=0 go install -ldflags \
  "-w -s -X main.build=${BUILD_NUMBER} -X main.commit=$(expr substr BUILD_COMMIT_SHORT 1 8)" \
  -a -tags netgo ./cmd/go-s3-backup

FROM scratch

COPY --from=builder /go/bin/go-s3-backup /go-s3-backup

ENTRYPOINT ["/go-s3-backup"]
