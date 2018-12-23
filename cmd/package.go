package cmd

import (
	"fmt"
	"os"

	"github.com/vinchauhan/goiib/util"

	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package the application by applying specific bar overrides",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := util.ApplyBarOverrides()
		if err != nil {
			fmt.Printf("Error running the bar override")
		}
		fmt.Println(result)
	},
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	RootCmd.AddCommand(packageCmd)
}
