ARG APP_IMAGE=codestation/go-s3-backup
FROM $APP_IMAGE AS binary

FROM postgres:10
LABEL maintainer="codestation <codestation404@gmail.com>"

RUN set -x \
	&& apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=binary /go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
