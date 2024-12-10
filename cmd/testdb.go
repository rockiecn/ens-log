package cmd

import (
	"github.com/rockiecn/ens-log/database"
	"github.com/urfave/cli/v2"
)

// read a name from db
var TestDBCmd = &cli.Command{
	Name:  "testdb",
	Usage: "test db",
	Action: func(cctx *cli.Context) error {
		database.TestDB()

		return nil
	},
}
