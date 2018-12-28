package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vinchauhan/goiib/config"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package the application by applying specific bar overrides",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := config.ApplyBarOverrides()
		if err != nil {
			fmt.Printf("Error running the bar override")
		}
		fmt.Println(result)
	},
}

//CreateDirIfNotExisting will create a directory passed if not exists
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
