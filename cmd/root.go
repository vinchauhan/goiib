package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "goiib",
	Short: "goiib is a build tool for IIB projects",
	Long:  `goiib clean/package/compile`,
}
