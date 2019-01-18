package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vinchauhan/goiib/config"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package the application by applying specific bar overrides",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Scanning for projects...")
		log.Infof("")

		_, err := config.ApplyBarOverrides()
		if err != nil {
			fmt.Printf("Error running the bar override")
		}
		log.Infof("-----------------------------------------------")
		log.Infof("BUILD SUCCESS")
		log.Infof("-----------------------------------------------")
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
