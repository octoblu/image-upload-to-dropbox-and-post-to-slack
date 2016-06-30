package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
			Name:   "example, e",
			EnvVar: "IMAGE_UPLOAD_TO_DROPBOX_AND_POST_TO_SLACK_EXAMPLE",
			Usage:  "Example string flag",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	example := getOpts(context)

	sigTerm := make(chan os.Signal)
	signal.Notify(sigTerm, syscall.SIGTERM)

	sigTermReceived := false

	go func() {
		<-sigTerm
		fmt.Println("SIGTERM received, waiting to exit")
		sigTermReceived = true
	}()

	for {
		if sigTermReceived {
			fmt.Println("I'll be back.")
			os.Exit(0)
		}

		debug("image-upload-to-dropbox-and-post-to-slack.loop: %v", example)
		time.Sleep(1 * time.Second)
	}
}

func getOpts(context *cli.Context) string {
	example := context.String("example")

	if example == "" {
		cli.ShowAppHelp(context)

		if example == "" {
			color.Red("  Missing required flag --example or IMAGE_UPLOAD_TO_DROPBOX_AND_POST_TO_SLACK_EXAMPLE")
		}
		os.Exit(1)
	}

	return example
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
