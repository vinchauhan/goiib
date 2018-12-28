package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinchauhan/goiib/config"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile IIB Application/Project and prepare a bar file",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := config.CompileProject()
		if err != nil {
			fmt.Printf("Could not compile the Project %v", err)
		}
		fmt.Println(result)
	},
}

func init() {
	RootCmd.AddCommand(compileCmd)
}
