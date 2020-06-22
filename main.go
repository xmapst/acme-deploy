package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"runtime"
)

var version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "acme"
	app.HelpName = "acme"
	app.Usage = "Let's Encrypt client written in Go"
	app.EnableBashCompletion = true
	app.Version = version
	app.Flags = createFlags()
	app.Commands = createCommands()
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("acme version %s %s/%s\n", c.App.Version, runtime.GOOS, runtime.GOARCH)
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
