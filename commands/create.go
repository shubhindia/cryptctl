package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	secretsv1alpha1 "github.com/shubhindia/encrypted-secrets/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create [flags] <cluster_name>",
	Short: "create encryptedSecrets manifest",
	Long:  "Create an EncryptedSecretinstance manifest file that will contain encrypted secret values",
	Args: func(_ *cobra.Command, args []string) error {

		if len(args) == 0 {

			os.Exit(1)
		}

		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {

		// get values from the flags
		fileName := args[0]
		provider := args[1]

		sampleEncryptedSecret := &secretsv1alpha1.EncryptedSecret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "encryptedsecret-sample",
				Namespace: "default",
			},
			TypeMeta: v1.TypeMeta{
				APIVersion: SecretApiVersion,
				Kind:       "EncryptedSecret",
			},
		}

		// set some sample labels
		sampleEncryptedSecret.SetLabels(map[string]string{
			"app.kubernetes.io/name":       "encryptedsecret",
			"app.kubernetes.io/part-of":    "encryted-secrets",
			"app.kubernetes.io/created-by": "encryted-secrets",
		})

		// set provider annotation
		sampleEncryptedSecret.SetAnnotations(map[string]string{
			"secrets.shubhindia.xyz/provider": provider,
		})

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
