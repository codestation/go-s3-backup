FROM golang:1.10-alpine as builder

COPY . $GOPATH/src/megpoid.xyz/go/gogs-s3-backup/
WORKDIR $GOPATH/src/megpoid.xyz/go/gogs-s3-backup/

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/gogs-s3-backup

FROM gogs/gogs:0.11.53

COPY --from=builder /go/bin/gogs-s3-backup /gogs-s3-backup

CMD ["/gogs-s3-backup"]
