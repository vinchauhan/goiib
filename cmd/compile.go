package cmd

import (
	"fmt"

	"github.com/vinchauhan/goiib/util"

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile IIB Application/Project and prepare a bar file",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := util.CompileProject()
		if err != nil {
			fmt.Printf("Error Compiling the code")
		}
		fmt.Println(result)
	},
}

func init() {
	RootCmd.AddCommand(compileCmd)
}
