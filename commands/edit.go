package commands

import (
	"fmt"
	"os"

	"github.com/shubhindia/hcictl/common"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	editutils "github.com/shubhindia/hcictl/commands/utils/edit"
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
			encryptedFile, err := os.ReadFile(fileName)
			if err != nil {
				return fmt.Errorf("error reading file %s", err.Error())
			}

			var encryptedSecret editutils.EncryptedSecret

			// unmarshal into EncryptedSecret
			err = yaml.Unmarshal(encryptedFile, &encryptedSecret)
			if err != nil {
				return fmt.Errorf("error unmarshaling file %s", err.Error())
			}

			// prepare decryptedSecret to be edited
			decryptedSecret := editutils.DecryptedSecret{
				ApiVersion: encryptedSecret.ApiVersion,
				Kind:       "DecryptedSecret",
				Metadata:   encryptedSecret.Metadata,
			}

			keyPhrase := os.Getenv("KEYPHRASE")
			if keyPhrase == "" {
				return fmt.Errorf("keyphrase not found")
			}

			decryptedData := make(map[string]string)

			// decrypt the data in encryptedSecrets
			for key, value := range encryptedSecret.Data {
				decryptedString := editutils.DecodeAndDecrypt(value, keyPhrase)

				decryptedData[key] = string(decryptedString)
			}

			decryptedSecret.Data = decryptedData

			fmt.Printf("%+v", decryptedSecret)

			return nil

		},
	}

	common.RegisterCommand(cliCmd)
}
