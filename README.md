# Go S3 Backup

This application can be used to make regular backups of various applications and restores from S3. It will perform daily backups unless configured to do otherwise.

## Supported apps

* PostgreSQL
* MySQL
* Gogs

The schedule function can also be used on restore if you need to test your backups regularly.

## Environment variables

* `S3_ENDPOINT`: url of the s3 region, for example `https://nyc3.digitaloceanspaces.com`.
* `S3_REGION`: region where the bucket is located, for example `us-east-1`.
* `S3_BUCKET`: name of the bucket, for example `backups`.
* `S3_PREFIX`: for example `private/files`.
* `S3_FORCE_PATH_STYLE`: set to `1` if you are using minio.
* `SCHEDULE`: specifies when to start a task. Defaults to `@daily` on backup, `none` on restore. Accepts cron format, like `0 0 * * * `. Set to `none` to disable and perform only one task.
* `SCHEDULE_RANDOM_DELAY`: maximum number of seconds (value choosen at random) to wait before starting a task. There is no random delay by default.
* `MAX_BACKUPS`: maximum number of backups to keep on S3. Only used with the `backup` command.
* `S3_KEY`: Restore directly from this S3 object instead of searching for the most recent one. Only used with the `restore` command.
* `DATABASE_HOST`: database host.
* `DATABASE_PORT`: database port.
* `DATABASE_USER`:  database user.
* `DATABASE_PASSWORD`:  database password.
