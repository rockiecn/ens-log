package cmd

import (
	"github.com/rockiecn/ens-log/common"
	"github.com/rockiecn/ens-log/utils"
	"github.com/urfave/cli/v2"
)

var SetLBCmd = &cli.Command{
	Name:  "setlocal",
	Usage: "set local block",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "block",
			Aliases: []string{"b"},
			Usage:   "block number to set",
		},
	},
	Action: func(cctx *cli.Context) error {
		b := cctx.String("b")
		i := utils.BytestoUint64([]byte(b))
		common.SetLocalBlock(i)
		return nil
	},
}
