FROM golang:1.11-alpine as builder

ARG BUILD_NUMBER=0
ARG BUILD_COMMIT_SHORT=unknown
ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /app
COPY . .

RUN go install -mod vendor -ldflags "-w -s \
  -X main.AppVersion=0.1.${BUILD_NUMBER:-0} \
  -X main.BuildCommit=${BUILD_COMMIT_SHORT} \
  -X \"main.BuildTime=$(date -u '+%Y-%m-%d %I:%M:%S %Z')\"" \
  -a ./cmd/go-s3-backup

FROM alpine:3.8
LABEL maintainer="codestation <codestation404@gmail.com>"

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
