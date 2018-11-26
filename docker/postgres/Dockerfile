ARG REPOSITORY=registry.megpoid.xyz
FROM $REPOSITORY/go-s3-backup AS app

FROM ubuntu:18.04
LABEL maintainer="codestation <codestation404@gmail.com>"

RUN set -ex; \
	apt-get update; \
	DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates wget gnupg2; \
        wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -; \
	DEBIAN_FRONTEND=noninteractive apt-get purge -y wget gnupg2; \
	rm -rf /var/lib/apt/lists/*

RUN set -ex; \
        echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" > /etc/apt/sources.list.d/pgdg.list; \
        apt-get update; \
        DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends postgresql-client-11; \
        rm -rf /var/lib/apt/lists/*

COPY --from=app /bin/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
