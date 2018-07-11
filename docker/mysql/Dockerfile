ARG APP_IMAGE=codestation/go-s3-backup
FROM $APP_IMAGE AS app

FROM alpine:3.8
LABEL maintainer="codestation <codestation404@gmail.com>"

RUN apk add --no-cache ca-certificates mysql-client
COPY --from=app /bin/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
