package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes/scheme"

	apis "github.com/shubhindia/cryptctl/apis"
)

var rootCmd = &cobra.Command{
	Use:           "cryptctl",
	Short:         "cryptctl is a command line tool for managing EncryptedSecrets",
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no command specified")
		}
		return nil
	},
}

func init() {
	_ = apis.AddToScheme(scheme.Scheme)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
}
