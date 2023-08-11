package cmd

import (
	"fmt"

	"github.com/rockiecn/ens-log/utils"
	"github.com/urfave/cli/v2"
)

// read a name from db
var ReadCmd = &cli.Command{
	Name:  "read",
	Usage: "test read a name from db",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Value:   "codingsh",
			Usage:   "name to read from db",
		},
	},
	Action: func(cctx *cli.Context) error {
		name := cctx.String("name")
		if name == "" {
			return (fmt.Errorf("name should not be blank"))
		}
		utils.Read(name)

		return nil
	},
}
