package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/image-upload-to-dropbox-and-post-to-slack/uploader"
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
			Name:   "content, c",
			EnvVar: "IUTDAPTS_CONTENT",
			Usage:  "Base64 encoded image to be uploaded to dropbox",
		},
		cli.StringFlag{
			Name:   "dropbox-access-token, d",
			EnvVar: "IUTDAPTS_DROPBOX_ACCESS_TOKEN",
			Usage:  "Dropbox Access Token",
		},
		cli.StringFlag{
			Name:   "dropbox-file-path, r",
			EnvVar: "IUTDAPTS_DROPBOX_FILE_PATH",
			Usage:  "Remote file path on Dropbox the image will be uploaded to",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	contentStrBase64, dropboxAccessToken, filePath := getOpts(context)

	dropbox := uploader.New(dropboxAccessToken)
	publicURL, err := dropbox.UploadBase64(filePath, contentStrBase64)
	fatalIfErr(err)

	fmt.Println("publicURL: ", publicURL)
}

func getOpts(context *cli.Context) (string, string, string) {
	contentStrBase64 := context.String("content")
	dropboxAccessToken := context.String("dropbox-access-token")
	dropboxFilePath := context.String("dropbox-file-path")

	if contentStrBase64 == "" || dropboxAccessToken == "" || dropboxFilePath == "" {
		cli.ShowAppHelp(context)

		if contentStrBase64 == "" {
			color.Red("  Missing required flag --content or IUTDAPTS_CONTENT")
		}
		if dropboxAccessToken == "" {
			color.Red("  Missing required flag --dropbox-access-token or IUTDAPTS_DROPBOX_ACCESS_TOKEN")
		}
		if dropboxFilePath == "" {
			color.Red("  Missing required flag --dropbox-file-path or IUTDAPTS_DROPBOX_FILE_PATH")
		}
		os.Exit(1)
	}

	return contentStrBase64, dropboxAccessToken, dropboxFilePath
}

func fatalIfErr(err error) {
	if err == nil {
		return
	}

	log.Fatalln(err.Error())
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
