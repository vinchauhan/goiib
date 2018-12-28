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

	// fmt.Println("Invoking main")
	// path := filepath.Join(filepath.Dir("."), "build.yaml")
	// //fmt.Println(path)

	// source, err := ioutil.ReadFile(path)
	// if err != nil {
	// 	panic(err)
	// }
	// err = yaml.Unmarshal(source, &config)
	// if err != nil {
	// 	panic(err)
	// }

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
