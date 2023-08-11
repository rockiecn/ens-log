package cmd

import (
	"fmt"
	"log"

	"github.com/rockiecn/ens-log/common"
	"github.com/urfave/cli/v2"
)

var GetLBCmd = &cli.Command{
	Name:  "getlocal",
	Usage: "get local block",
	Action: func(cctx *cli.Context) error {
		l, err := common.GetLocalBlock()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("local block:", l)
		return nil
	},
}
