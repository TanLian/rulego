package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/tanlian/rulego/program"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run scripts with the .rg extension file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("egg: rg run xxx.rg")
			return
		}

		filePath := args[0]
		f, err := os.Open(filePath)
		if err != nil {
			fmt.Println("run err: ", err)
			return
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			fmt.Println("run err: ", err)
			return
		}
		program.New().Run(string(data))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
