package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/garden"
	gclient "code.cloudfoundry.org/garden/client"
	gconn "code.cloudfoundry.org/garden/client/connection"
	"github.com/urfave/cli"
	"github.com/williammartin/garden-curator/blueprint"
	yaml "gopkg.in/yaml.v2"
)

var GrowCommand = cli.Command{
	Name: "grow",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "blueprint, b",
			Usage: "The filepath to a blueprint.yml.",
		},
	},
	Action: func(ctx *cli.Context) error {
		blueprintFile := ctx.String("blueprint")
		if blueprintFile == "" {
			blueprintFile = "blueprint.yml"
		}

		// unmarshal yml
		bytes, err := ioutil.ReadFile(blueprintFile)
		if err != nil {
			return err
		}

		blueprint := &blueprint.Blueprint{}
		err = yaml.Unmarshal(bytes, blueprint)
		if err != nil {
			return err
		}

		client := gclient.New(gconn.New("tcp", "10.244.0.2:7777"))
		for _, handle := range blueprint.Containers {
			fmt.Printf("growing '%s'...\n", handle)
			_, err = client.Create(garden.ContainerSpec{Handle: handle})
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "garden-curator"
	app.Commands = []cli.Command{
		GrowCommand,
	}

	app.Run(os.Args)
}
