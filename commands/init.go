package commands

import (
	"github.com/shubhindia/cryptctl/commands/utils"
	"github.com/spf13/cobra"
)

func init() {

	initCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "namespace to use (required)")
	initCmd.MarkFlagRequired("namespace")
	initCmd.Flags().StringVarP(&Provider, "provider", "p", "", "provider to use (required)")
	initCmd.MarkFlagRequired("provider")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [flags]",
	Short: "init",
	Long:  "Init initializes the encrypted-secrets CLI",

	RunE: func(_ *cobra.Command, args []string) error {

		switch Provider {
		case "k8s":
			return utils.InitK8s(Namespace)
		}

		return nil

	},
}
