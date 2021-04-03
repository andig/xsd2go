package cmd

import (
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func init() {
	app.Name = "XSD2Go"
	app.Usage = "Automatically generate golang xml parser based on XSD"
}

func Execute() error {
	return app.Run(os.Args)
}
