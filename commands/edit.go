package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/shubhindia/hcictl/commands/utils/edit"
	"github.com/shubhindia/hcictl/common"
	"github.com/urfave/cli/v2"
)

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

			fileName := ctx.Args().First()

			// read the file
			var inStream io.Reader
			inFile, err := os.Open(fileName)
			if err != nil {
				if os.IsNotExist(err) {
					return errors.Wrapf(err, "error reading input file %s", fileName)

				} else {
					return errors.Wrapf(err, "error reading input file %s", fileName)
				}
			} else {
				defer inFile.Close()
				inStream = inFile
			}

			// Parse the input file to objects.
			inManifest, err := edit.NewManifest(inStream)
			if err != nil {
				return errors.Wrap(err, "error decoding input YAML")
			}

			fmt.Printf("%+v", inManifest)

			return nil
		},
	}

	common.RegisterCommand(cliCmd)
}
