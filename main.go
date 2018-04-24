package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "garden-curator"
	app.Commands = []cli.Command{}

	app.Run(os.Args)
}
