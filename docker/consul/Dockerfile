ARG REPOSITORY=registry.megpoid.xyz
FROM $REPOSITORY/go-s3-backup AS app

FROM consul:1.2.1
LABEL maintainer="codestation <codestation404@gmail.com>"

COPY --from=app /bin/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
