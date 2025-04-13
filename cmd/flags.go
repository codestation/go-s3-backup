package cmd

import "github.com/spf13/pflag"

func LoadDefaultFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.Int("schedule-random-delay", 1, "Schedule random delay")
	fs.String("save-dir", "/tmp/go-s3-backup", "Directory to save/read backups")
	return fs
}

func LoadBackupFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.String("schedule", "@daily", "Cron schedule")
	fs.Int("max-backups", 5, "Max backups to keep (0 to disable the feature)")
	return fs
}

func LoadRestoreFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.String("schedule", "none", "Cron schedule")
	fs.String("restore-file", "", "Restore from this file instead of searching for the most recent")
	fs.String("restore-prefix", "", "Name prefix to filter when restoring the backup")
	return fs
}

func LoadDatabaseFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.String("database-host", "", "Database host")
	fs.String("database-port", "", "Database port")
	fs.String("database-name", "", "Database name")
	fs.String("database-user", "", "Database user")
	fs.String("database-password", "", "Database password")
	fs.String("database-password-file", "", "Database password file")
	fs.String("database-filename-prefix", "", "Database filename prefix")
	fs.String("database-name-as-prefix", "", "Database name as prefix")
	fs.String("database-options", "", "Extra options to pass to database service")
	fs.String("database-compress", "", "Compress sql with gzip")
	fs.String("database-ignore-exit-code", "", "Ignore restore process exit code")
	return fs
}

func LoadMySQLFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.Bool("mysql-split-databases", false, "Make individual backups instead of a single one")
	fs.StringSlice("mysql-exclude-databases", nil, "Make backup of databases except the ones that matches the pattern")
	return fs
}

func LoadPostgresFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.String("postgres-binary-path", "", "Directory where postgres binaries are located")
	fs.String("postgres-version", "17", "Postgres version for the pg_dump/pg_restore/psql tools")
	fs.Bool("postgres-custom-format", false, "Use custom format (always compressed), ignored when database name is not set")
	fs.Bool("postgres-drop", false, "Drop the database before restoring it")
	fs.String("postgres-owner", "", "Change owner on database restore")
	fs.StringSlice("postgres-exclude-databases", nil, "Make backup of databases except the ones that matches the pattern")
	fs.Bool("postgres-backup-per-user", false, "Make backups for all databases separated per user")
	fs.StringSlice("postgres-backup-users", nil, "Make backups for databases matching these users")
	fs.StringSlice("postgres-backup-exclude-users", []string{"postgres"}, "Make backups for databases excluding these users")
	fs.Bool("postgres-backup-per-schema", false, "Make backups separated per schema")
	fs.StringSlice("postgres-backup-schemas", nil, "Make backups matching these schemas")
	fs.StringSlice("postgres-backup-exclude-schemas", []string{"information_schema", "pg_toast", "pg_catalog"}, "Make backup excluding these schemas")
	return fs
}

func LoadTarballFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.String("tarball-name-prefix", "", "Backup file prefix")
	fs.String("tarball-path-source", "", "Path to backup/restore")
	fs.Bool("tarball-compress", false, "Compress tarball with gzip")
	fs.String("tarball-path-prefix", "", "Backup path prefix")
	fs.Bool("tarball-backup-per-dir", false, "Backup each folder individually")
	fs.StringSlice("tarball-backup-dirs", nil, "Backup each folder individually")
	fs.StringSlice("tarball-backup-exclude-dirs", nil, "Make backups for directories excluding these dirs")
	return fs
}

func LoadS3Flags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.String("s3-endpoint", "", "S3 endpoint")
	fs.String("s3-region", "", "S3 region")
	fs.String("s3-bucket", "", "S3 bucket")
	fs.String("s3-prefix", "", "S3 prefix")
	fs.Bool("s3-force-path-style", false, "S3 force path style (needed for minio)")
	fs.Bool("s3-keep-file", false, "Keep local file after successful upload")
	return fs
}
