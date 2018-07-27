# Go S3 Backup

This application can be used to make regular backups of various applications and restores from S3. It will perform daily backups unless configured to do otherwise.

## Supported services

* PostgreSQL
* MySQL
* Gogs
* Tarball
* Consul

## Supported stores
* S3
* Filesystem (local)

The schedule function can also be used on restore if you need to test your backups regularly.

## Environment variables

### Global configuration
* `CONFIG`: load config from a yaml file

### Backup/restore configuration
* `SAVE_DIR`: directory to store the temporal backup after creating/retrieving it.`
* `SCHEDULE_RANDOM_DELAY`: maximum number of seconds (value choosen at random) to wait before starting a task. There is no random delay by default.
* `SCHEDULE`: specifies when to start a task. Defaults to `@daily` on backup, `none` on restore. Accepts cron format, like `0 0 * * * `. Set to `none` to disable and perform only one task.

### Backup-related configuration
* `MAX_BACKUPS`: maximum number of backups to keep on the store.

### Restore related configuration
* `RESTORE_FILE`: Restore directly from this filename instead of searching for the most recent one. Only used with the `restore` command.

### Gogs configuration
* `GOGS_CONFIG`: custom location of the gogs config file.
* `GOGS_DATA`: location of the Gogs data directory.

### Database common config
* `DATABASE_HOST`: database host.
* `DATABASE_PORT`: database port.
* `DATABASE_NAME`: database name.
* `DATABASE_USER`:  database user.
* `DATABASE_PASSWORD`:  database password.
* `DATABASE_PASSWORD_FILE`:  database password file, has precendnce over `DATABASE_PASSWORD`
* `DATABASE_OPTIONS`:  custom options to pass to the backup/restore application.
* `DATABASE_COMPRESS`: compress the sql file with gzip.
* `DATABASE_IGNORE_EXIT_CODE`: ignore is the restore operation returns a non-zero exit code.

### Postgres configuration
* `POSTGRES_CUSTOM_FORMAT`: use custom dump format instead of plain text backups.

### Tarball configuration
* `TARBALL_PATH_SOURCE`: directory to backup/restore.
* `TARBALL_NAME_PREFIX`: name prefix of the created tarball. If unset it will use the backup directory name.
* `TARBALL_COMPRESS`: compress the tarball with gzip.

### S3 configuration
* `S3_ENDPOINT`: url of the 33 endpoint, for example `https://nyc3.digitaloceanspaces.com`.
* `S3_REGION`: region where the bucket is located, for example `us-east-1`.
* `S3_BUCKET`: name of the bucket, for example `backups`.
* `S3_PREFIX`: for example `private/files`.
* `S3_FORCE_PATH_STYLE`: set to `1` if you are using minio.
* `S3_KEEP_FILE`: keep file on the local filesystem after uploading it to S3.

The credentials are passed using the standard variables:
* `AWS_ACCESS_KEY_ID`: AWS access key. `AWS_ACCESS_KEY` can also be used.
* `AWS_SECRET_ACCESS_KEY`: AWS secret key. `AWS_SECRET_KEY` can also be used.
* `AWS_SESSION_TOKEN`: AWS session token. Optional, will be used if present.

The credentials can also be stored in a file. The location of the file can be set in `AWS_SHARED_CREDENTIALS_FILE`, else it will use `~/.aws/credentials` The format of the file is the following:

``` ini
[default]
aws_access_key_id = YOUR_AWS_ACCESS_KEY_ID
aws_secret_access_key = YOUR_AWS_SECRET_ACCESS_KEY
```
