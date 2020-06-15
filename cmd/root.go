package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yyh-gl/go-ec2/internal"
)

var (
	configPath string

	rootCmd = &cobra.Command{
		Use:     "go-ec2",
		Version: "0.2.1",
		Short:   "EC2 Manager",
		Long:    "Simple EC2 Manager made by Go.",
	}

	stateCmd = &cobra.Command{
		Use:   "state",
		Short: "Print the states of all instances",
		Long:  "Print the states of all instances",
		Run: func(cmd *cobra.Command, args []string) {
			m := internal.NewManger(configPath)
			if err := m.PrintAllState(context.Background()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Print the states of all instances",
		Long:  "Print the states of all instances",
		Run: func(cmd *cobra.Command, args []string) {
			m := internal.NewManger(configPath)
			if err := m.StopAllInstances(context.Background()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stateCmd.Flags().StringVarP(&configPath, "cfgPath", "c", homeDir+"/.go-ec2.yml", "Path to config file")
	stopCmd.Flags().StringVarP(&configPath, "cfgPath", "c", homeDir+"/.go-ec2.yml", "Path to config file")
	rootCmd.AddCommand(stateCmd)
	rootCmd.AddCommand(stopCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
