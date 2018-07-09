ARG APP_IMAGE=codestation/go-s3-backup
FROM $APP_IMAGE AS binary

FROM postgres:10
LABEL maintainer="codestation <codestation404@gmail.com>"

COPY --from=binary /go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
