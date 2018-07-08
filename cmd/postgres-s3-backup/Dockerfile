FROM golang:1.10-alpine as builder

COPY . $GOPATH/src/megpoid.xyz/go/go-s3-backup/
WORKDIR $GOPATH/src/megpoid.xyz/go/go-s3-backup/

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/postgres-s3-backup

FROM postgres:10

COPY --from=builder /go/bin/postgres-s3-backup /bin/postgres-s3-backup

CMD ["/bin/postgres-s3-backup"]
