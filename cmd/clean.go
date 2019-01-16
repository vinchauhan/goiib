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

		targetDirPath := filepath.Join(buildConfig.Project.Profiles.Profile.Properties.Workspace, buildConfig.Project.ArtifactID, "target")
		err := CleanProject(targetDirPath)
		if err != nil {
			fmt.Printf("% v : Error in cleaning the target dir %s", err, targetDirPath)
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

	path := filepath.Join(filepath.Dir("."), "build.yaml")

	source, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("goiib error : %v", err)
	}
	err = yaml.Unmarshal(source, &buildConfig)
	if err != nil {
		log.Fatalf("goiib error : Could not Unmarshal the build %v", err)
	}

}

// CleanProject will clean the target directory
func CleanProject(targetPath string) error {
	log.Info("---goiib clean command @ ", buildConfig.Project.ArtifactID, " ---")
	log.Infof("Deleting %s", targetPath)
	err := os.RemoveAll(targetPath)
	if err != nil {
		return fmt.Errorf("Error Occured: %v", err)
	}
	log.Infof("Deleted %s", targetPath)
	return nil
}
