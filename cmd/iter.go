package cmd

import (
	"log"

	"github.com/rockiecn/ens-log/common"
	"github.com/urfave/cli/v2"
)

var IterCmd = &cli.Command{
	Name:  "iter",
	Usage: "iterate k-v",
	Action: func(cctx *cli.Context) error {
		err := common.Iter()
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}
