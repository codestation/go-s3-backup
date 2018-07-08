FROM golang:1.10-alpine as builder

COPY . $GOPATH/src/megpoid.xyz/go/go-s3-backup/
WORKDIR $GOPATH/src/megpoid.xyz/go/go-s3-backup/

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/mysql-s3-backup

FROM mysql:5.7

COPY --from=builder /go/bin/mysql-s3-backup /bin/mysql-s3-backup

CMD ["/bin/mysql-s3-backup"]
