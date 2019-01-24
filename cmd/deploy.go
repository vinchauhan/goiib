package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vinchauhan/goiib/config"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an IIB Application/Library and store binary on Sonatype Nexus Repository",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := config.DeployProject(buildFilePath)
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
	RootCmd.AddCommand(deployCmd)
}
