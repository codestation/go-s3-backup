ARG APP_IMAGE=codestation/go-s3-backup
FROM $APP_IMAGE AS app

FROM mysql:5.7
LABEL maintainer="codestation <codestation404@gmail.com>"

COPY --from=app /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=app /bin/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
