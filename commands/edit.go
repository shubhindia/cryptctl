package commands

import (
	"fmt"
	"regexp"

	"github.com/shubhindia/hcictl/common"
	"github.com/urfave/cli/v2"
)

var whitespaceRegexp *regexp.Regexp

func init() {

	cliCmd := cli.Command{
		Name:  "edit",
		Usage: "edit encryptedSecrets manifest",
		Before: func(ctx *cli.Context) error {
			if ctx.Args().First() == "" {
				return fmt.Errorf("hcictl edit expectes a file to edit")
			}

			if ctx.Args().Len() > 1 {
				return fmt.Errorf("too many arguments")
			}

			return nil
		},
		Action: func(ctx *cli.Context) error {

			return nil

		},
	}

	common.RegisterCommand(cliCmd)
}
