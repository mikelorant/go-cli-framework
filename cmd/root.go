package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	_configEnvPrefix = "TEST5"
	_configFilename  = "test5.yaml"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test5",
		Short: "A brief description of your application",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(NewSSHCmd())

	cmd.PersistentFlags().String("config", "", "config file (default is $HOME/test5.yaml)")
	cmd.PersistentFlags().Bool("verbose", false, "enable verbose output")

	cmd.PersistentFlags().SetAnnotation("config", "bindWithKey", []string{"config"})
	cmd.PersistentFlags().SetAnnotation("verbose", "bindWithKey", []string{"debug"})

	return cmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "error: %v", err)
		os.Exit(1)
	}
}
