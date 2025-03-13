package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "no-fines",
	Short: "no-fines app",
	Long:  "long info about no-fines app",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use 'serve' for run service.")
		fmt.Println("use 'migrate' for apply migration in database.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
