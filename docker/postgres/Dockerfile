ARG REPOSITORY=registry.megpoid.xyz
FROM $REPOSITORY/go-s3-backup AS app

FROM ubuntu:18.04
LABEL maintainer="codestation <codestation404@gmail.com>"

RUN set -ex; \
        if ! command -v gpg > /dev/null; then \
                apt-get update; \
                apt-get install -y --no-install-recommends \
                        gnupg \
                        dirmngr \
                ; \
                rm -rf /var/lib/apt/lists/*; \
        fi

RUN set -ex; \
        key='B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8'; \
        export GNUPGHOME="$(mktemp -d)"; \
        gpg --batch --keyserver ha.pool.sks-keyservers.net --recv-keys "$key"; \
        gpg --batch --export "$key" > /etc/apt/trusted.gpg.d/postgres.gpg; \
        command -v gpgconf > /dev/null && gpgconf --kill all; \
        rm -rf "$GNUPGHOME"; \
        apt-key list

RUN set -ex; \
        echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" > /etc/apt/sources.list.d/pgdg.list; \
        apt-get update; \
        DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates postgresql-client-11; \
        rm -rf /var/lib/apt/lists/*

COPY --from=app /bin/go-s3-backup /bin/go-s3-backup

ENTRYPOINT ["/bin/go-s3-backup"]
