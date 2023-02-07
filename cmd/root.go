package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "wren task syncer",
	Long:  "Sync wren tasks from polcode link app",
}

func init() {
	rootCmd.AddCommand(SyncCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
