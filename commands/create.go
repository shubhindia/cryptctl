package commands

import (
	"fmt"
	"os"

	"github.com/shubhindia/cryptctl/common"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	editutils "github.com/shubhindia/cryptctl/commands/utils/edit"
)

const (
	apiVersion = "secrets.shubhindia.xyz/v1alpha1"
	kind       = "EncryptedSecret"
)

func init() {

	createCmd := cli.Command{
		Name:  "create",
		Usage: "create encryptedSecrets manifest",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "filename",
				Aliases: []string{
					"f",
				},
				Required: true,
			},
			&cli.StringFlag{
				Name: "provider",
				Aliases: []string{
					"p",
				},
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {

			// get values from the flags
			fileName := ctx.String("filename")
			provider := ctx.String("provider")

			sampleEncryptedSecret := editutils.EncryptedSecret{
				ApiVersion: apiVersion,
				Kind:       kind,
			}

			metadata := make(map[string]interface{})
			metadata["name"] = "encryptedsecret-sample"
			metadata["namespace"] = "default"

			// add some default labels
			labelsMap := map[string]string{
				"app.kubernetes.io/created-by": "encryted-secrets",
				"app.kubernetes.io/instance":   "encryptedsecret-sample",
				"app.kubernetes.io/name":       "encryptedsecret",
				"app.kubernetes.io/part-of":    "encryted-secrets",
				"app.kubernetes.io/provider":   provider,
			}
			metadata["labels"] = labelsMap

			sampleEncryptedSecret.Metadata = metadata

			// write the contents to yaml
			newEncrypted, err := yaml.Marshal(&sampleEncryptedSecret)
			if err != nil {
				return fmt.Errorf("error marshaling encryptedSecret %s", err.Error())
			}

			err = os.WriteFile(fileName, newEncrypted, 0600)
			if err != nil {
				return fmt.Errorf("error writing EncryptedSecret %s", err)
			}

			return nil

		},
	}

	common.RegisterCommand(createCmd)
}
