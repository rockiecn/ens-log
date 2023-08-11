package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rockiecn/ens-log/cmd"
)

/*
Filter nameRegister event logs for ethRegisterController of ENS.
*/
func main() {
	commands := []*cli.Command{
		cmd.UpdateCmd,
		cmd.ReadCmd,
		cmd.TestCmd,
		cmd.SetLBCmd,
		cmd.GetLBCmd,
		cmd.IterCmd,
	}

	app := &cli.App{
		Name:                 "ens-log",
		Usage:                "ens-log update",
		Version:              "v1",
		EnableBashCompletion: true,
		Commands:             commands,
	}

	app.Setup()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
