ARG REPOSITORY=registry.megpoid.xyz
FROM $REPOSITORY/go-s3-backup AS app

FROM alpine:3.9
LABEL maintainer="codestation <codestation404@gmail.com>"

RUN apk add --no-cache ca-certificates mysql-client
COPY --from=app /bin/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
