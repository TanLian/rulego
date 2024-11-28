package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tanlian/rulego/repl"
)

var rootCmd = &cobra.Command{
	Use: "rg",
	Run: func(cmd *cobra.Command, args []string) {
		repl.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
