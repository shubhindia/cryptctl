package commands

import (
	"fmt"
	"os"
	"regexp"

	secretsv1alpha1 "github.com/shubhindia/encrypted-secrets/api/v1alpha1"
	"github.com/shubhindia/encrypted-secrets/pkg/providers"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {

	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate [args]",
	Short: "validate",
	Long:  "Validate the encrypted-secrets yaml",
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return fmt.Errorf("filename is required")
		}
		return nil

	},

	RunE: func(_ *cobra.Command, args []string) error {

		// ToDo: commanize this code
		fileName := args[0]
		encryptedFile, err := os.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("error reading file %s", err.Error())
		}

		codecs := serializer.NewCodecFactory(scheme.Scheme, serializer.EnableStrict)
		obj, _, err := codecs.UniversalDeserializer().Decode(encryptedFile, &schema.GroupVersionKind{
			Group:   secretsv1alpha1.GroupVersion.Group,
			Version: SecretApiVersion,
			Kind:    "EncryptedSecret",
		}, nil)
		if err != nil {
			if ok, _ := regexp.MatchString("no kind(.*)is registered for version", err.Error()); ok {
				panic("no kind(.*)is registered for version")
			}
			panic(err)
		}

		// convert the runtimeObj to encryptedSecret object
		encryptedSecret, ok := obj.(*secretsv1alpha1.EncryptedSecret)
		if !ok {
			// should never happen
			panic("failed to convert runtimeObject to encryptedSecret")
		}

		// get the decryptedSecret
		decryptedObj, err := providers.DecodeAndDecrypt(encryptedSecret)
		if err != nil {
			return err
		}

		if len(encryptedSecret.Data) != len(decryptedObj.Data) {
			return fmt.Errorf("validation failed")
		}
		fmt.Println("validation successful")
		return nil

	},
}
