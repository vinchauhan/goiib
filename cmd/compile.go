package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vinchauhan/goiib/config"
)

var buildFile = "build.yaml"

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile IIB Application/Project and prepare a bar file",
	Run: func(cmd *cobra.Command, args []string) {

		log.Infof("Scanning for projects...")
		log.Infof("")

		_, err := config.CompileProject(buildFile)
		if err != nil {
			log.Errorf("goiib error : %v", err)
			log.Infof("-----------------------------------------------")
			log.Infof("BUILD FAILED")
			log.Infof("-----------------------------------------------")
		} else {
			log.Infof("-----------------------------------------------")
			log.Infof("BUILD SUCCESS")
			log.Infof("-----------------------------------------------")
		}

	},
}

func init() {
	RootCmd.AddCommand(compileCmd)
}
