package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinchauhan/goiib/config"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an IIB Application/Library and store binary on Sonatype Nexus Repository",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := config.DeployProject()
		if err != nil {
			fmt.Printf("Error Deploying the code")
		}
		fmt.Println(result)
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
