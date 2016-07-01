package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("image-upload-to-dropbox-and-post-to-slack:main")

func main() {
	app := cli.NewApp()
	app.Name = "image-upload-to-dropbox-and-post-to-slack"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "dropbox-access-token, d",
			EnvVar: "IUTDAPTS_DROPBOX_ACCESS_TOKEN",
			Usage:  "Dropbox Access Token",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	dropboxAccessToken := getOpts(context)
	debug(dropboxAccessToken)
}

func getOpts(context *cli.Context) string {
	dropboxAccessToken := context.String("dropbox-access-token")

	if dropboxAccessToken == "" {
		cli.ShowAppHelp(context)

		if dropboxAccessToken == "" {
			color.Red("  Missing required flag --dropbox-access-token or IUTDAPTS_DROPBOX_ACCESS_TOKEN")
		}
		os.Exit(1)
	}

	return dropboxAccessToken
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
