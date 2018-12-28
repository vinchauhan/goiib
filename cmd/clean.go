package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//var config = util.AutoGenerated{}

var dirToClean = "target"

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the target directory just like mvn clean",
	Run: func(cmd *cobra.Command, args []string) {

		err := CleanProject(dirToClean)
		if err != nil {
			fmt.Printf("Error deleting the %s directory %v", dirToClean, err)
		}
		log.Infof("BUILD SUCCESS : %s directory deleted", dirToClean)
	},
}

func init() {

	//Initialize Logger
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	RootCmd.AddCommand(cleanCmd)

	//path := filepath.Join(filepath.Dir("."), "build.yaml")
	//fmt.Println(path)

	// source, err := ioutil.ReadFile(path)
	// if err != nil {
	// 	panic(err)
	// }
	// err = yaml.Unmarshal(source, &config)
	// if err != nil {
	// 	panic(err)
	// }

}

//CleanProject will clean the target dir in the current dir
func CleanProject(targetDir string) error {
	err := os.RemoveAll(targetDir)
	if err != nil {
		return err
	}
	return nil
}
