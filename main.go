package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time
var backupPrefix = "gogs"
var appPath = "/app/gogs/gogs"

func main() {
	app := cli.NewApp()
	app.Usage = "drone-stack plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "s3 endpoint",
			EnvVar: "S3_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "region",
			Usage:  "s3 region",
			EnvVar: "S3_REGION",
		},
		cli.StringFlag{
			Name:   "bucket",
			Usage:  "s3 bucket",
			EnvVar: "S3_BUCKET",
		},
		cli.StringFlag{
			Name:   "prefix",
			Usage:  "s3 prefix",
			EnvVar: "S3_PREFIX",
		},
		cli.BoolFlag{
			Name:   "force-path-style",
			Usage:  "s3 force path style (needed for minio)",
			EnvVar: "S3_FORCE_PATH_STYLE",
		},
		cli.StringFlag{
			Name:   "schedule",
			Usage:  "cron schedule",
			Value:  "@daily",
			EnvVar: "CRON_SCHEDULE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
