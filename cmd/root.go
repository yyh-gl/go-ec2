package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configPath string

	rootCmd = &cobra.Command{
		Use:     "go-ec2",
		Version: "0.1.0",
		Short:   "EC2 Manager",
		Long:    "Simple EC2 Manager made by Go.",
		Run:     func(cmd *cobra.Command, args []string) {},
	}
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.Flags().StringVarP(&configPath, "cfgPath", "c", homeDir+"/.go-ec2.yml", "Path to config file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
