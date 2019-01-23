package main

import (
	"fmt"
	"os"

	"github.com/vinchauhan/goiib/cmd"
)

const (
	buildFile = "build.yaml"
)

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
