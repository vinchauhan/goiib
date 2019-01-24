package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vinchauhan/goiib/config"
	yaml "gopkg.in/yaml.v2"
)

var buildConfig = config.BuildConfig{}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the target directory just like mvn clean",
	Run: func(cmd *cobra.Command, args []string) {

		err := CleanProject(buildFilePath)
		if err != nil {
			fmt.Printf("% v : Error in cleaning the target dir", err)
		}
		log.Info("--------------------------------------------------")
		log.Info("BUILD SUCCESS")
		log.Info("--------------------------------------------------")
	},
}

func init() {

	//Initialize Logger
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	RootCmd.AddCommand(cleanCmd)
}

// CleanProject will clean the target directory
func CleanProject(buildFilePath string) error {
	log.Info("---goiib clean command @ ", buildConfig.Project.ArtifactID, " ---")

	source, err := ioutil.ReadFile(buildFilePath)
	if err != nil {
		return fmt.Errorf("goiib error : Error occured in reading the file : %s %v", buildFilePath, err)
	}
	err = yaml.Unmarshal(source, &buildConfig)
	if err != nil {
		return fmt.Errorf("goiib error : Could not Unmarshal the build %v", err)
	}

	targetDirPath := filepath.Join(buildConfig.Project.Profiles.Profile.Properties.Workspace, buildConfig.Project.ArtifactID, "target")

	log.Infof("Deleting %s", targetDirPath)
	err = os.RemoveAll(targetDirPath)
	if err != nil {
		return fmt.Errorf("Error Occured: %v", err)
	}
	log.Infof("Deleted %s", targetDirPath)
	return nil
}
