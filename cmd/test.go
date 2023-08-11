package cmd

import (
	"github.com/rockiecn/ens-log/utils"
	"github.com/urfave/cli/v2"
)

// test db and encoder
var TestCmd = &cli.Command{
	Name:  "test",
	Usage: "test db and encoder",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		utils.Test()
		return nil
	},
}
