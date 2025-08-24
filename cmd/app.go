package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "app",
		Short: "application",
		Long:  "application",
		Run:   func(_ *cobra.Command, _ []string) {},
	}
	rootCmd.AddCommand(serviceCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	return err
}
