package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

//CompileProject -cleanBuild the IIB project and creates a BAR file in the /target directory
func CompileProject() ([]byte, error) {

	config := BuildConfig{}
	path := filepath.Join(filepath.Dir("."), "build.yaml")
	//fmt.Println(path)

	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return nil, fmt.Errorf("Could not Unmarshal the build.yaml %v", err)
	}
	fmt.Printf("Value: %#v\n", config.Project.GroupID)

	var mqsiCreateBar = &MqsiCommand{
		mqsi:                 "mqsicreatebar",
		mqsiPath:             filepath.Join("opt", "ibm", "iib-10.0.0.10", "server", "bin", "mqsiprofile"),
		dataOption:           "-data",
		workspace:            config.Project.Profiles.Profile.Properties.Workspace,
		artifactID:           config.Project.ArtifactID,
		version:              config.Project.Version,
		barfileOption:        "-b",
		applicationOption:    "-a",
		deployAsSourceOption: "-deployAsSource",
		traceOption:          "-trace",
		verboseOption:        "-v",
		tracePath:            filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, "target"),
		traceFile:            "createbar.txt",
		cleanBuildOption:     "-cleanBuild",
		barfileName:          filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, "target", config.Project.ArtifactID+"-"+config.Project.Version+".bar"),
	}

	mqsiCreateBarCmdString := fmt.Sprintln(mqsiCreateBar.mqsi,
		mqsiCreateBar.dataOption,
		mqsiCreateBar.workspace,
		mqsiCreateBar.barfileOption,
		mqsiCreateBar.barfileName,
		mqsiCreateBar.applicationOption,
		mqsiCreateBar.artifactID,
		mqsiCreateBar.cleanBuildOption,
		mqsiCreateBar.deployAsSourceOption,
		mqsiCreateBar.traceOption,
		mqsiCreateBar.verboseOption,
		mqsiCreateBar.tracePath,
	)
	fmt.Printf("mqsiCreateBar command is : %s", mqsiCreateBarCmdString)

	// fmt.Println(mqsiCreateBar.mqsi,
	// 	mqsiCreateBar.dataOption,
	// 	mqsiCreateBar.workspace,
	// 	mqsiCreateBar.barfileOption,
	// 	mqsiCreateBar.barfileName,
	// 	mqsiCreateBar.applicationOption,
	// 	mqsiCreateBar.artifactID,
	// 	mqsiCreateBar.cleanBuildOption,
	// 	mqsiCreateBar.deployAsSourceOption,
	// 	mqsiCreateBar.traceOption,
	// 	mqsiCreateBar.verboseOption,
	// 	mqsiCreateBar.tracePath,
	// )

	cmd := exec.Command("/bin/bash", "-c", mqsiCreateBarCmdString)
	//cmd := exec.Command("/bin/bash", "-C", "ls -l;wc -l")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Errorf("Couldnt not execute the program %v", err)
	}

	// cmd := exec.Command(
	// 	mqsiCreateBar.mqsi,
	// 	mqsiCreateBar.dataOption,
	// 	mqsiCreateBar.workspace,
	// 	mqsiCreateBar.barfileOption,
	// 	mqsiCreateBar.barfileName,
	// 	mqsiCreateBar.applicationOption,
	// 	mqsiCreateBar.artifactID,
	// 	mqsiCreateBar.cleanBuildOption,
	// 	mqsiCreateBar.deployAsSourceOption,
	// 	mqsiCreateBar.traceOption,
	// 	mqsiCreateBar.verboseOption,
	// 	mqsiCreateBar.tracePath+mqsiCreateBar.traceFile)

	// out, err := cmd.Output()
	// if err != nil {
	// 	//return nil, err
	// 	panic(err)
	// }

	//Create defalt.properties file in the /target directory to use for creating the override file.

	//start reading the createbar.txt for more verbose output
	defaultPropCmd := exec.Command("mqsireadbar", "-b", mqsiCreateBar.barfileName, "-r")

	defaultout, defaulterr := defaultPropCmd.Output()
	if defaulterr != nil {
		return nil, fmt.Errorf("Could execute the command mqsireadbar and failed with %v", err)
	}

	err = ioutil.WriteFile(mqsiCreateBar.tracePath+"default.properties", []byte(defaultout), 755)

	if err != nil {
		log.Fatalf("mqsireadbar filed with %s\n", err)
	}
	//fmt.Printf("%s", out)
	return nil, nil

}

//This function will Deploy the bar file on the broker and push the binary to Nexus
func DeployProject() (string, error) {

	config := BuildConfig{}

	path := filepath.Join(filepath.Dir("."), "build.yaml")
	//fmt.Println(path)

	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}

	extraParams := map[string]string{
		"raw.directory":       "com/aa/esoa/iib",
		"raw.asset1.filename": config.Project.ArtifactID + "-" + config.Project.Version + ".bar",
	}

	fileName := config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\target\\iib-overrides\\" + config.Project.ArtifactID + "-" + config.Project.Version + "-dev" + ".bar"

	request, err := newfileUploadRequest("http://localhost:4444/service/rest/beta/components?repository=esoa-releases", extraParams, "raw.asset1", fileName)
	if err != nil {
		log.Fatal(err)
	}

	// request.SetBasicAuth("admin", "admin123")
	// //req.Header.Set("Accept-Encoding", "multipart/form-data")
	// request.Header.Set("Content-Type", "multipart/form-data")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.Status, nil
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth("admin", "admin123")
	return req, err
}

//ApplyBarOverrides will apply the bar override and create the overriden bar file in target/iib-overrides
func ApplyBarOverrides() ([]byte, error) {

	config := BuildConfig{}

	path := filepath.Join(filepath.Dir("."), "build.yaml")
	//fmt.Println(path)

	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Value: %#v\n", config.Project.GroupID)

	var mqsiBarOverride = &MqsiCommand{
		mqsi:                 "mqsiapplybaroverride",
		dataOption:           "-data",
		workspace:            config.Project.Profiles.Profile.Properties.Workspace,
		artifactID:           config.Project.ArtifactID,
		version:              config.Project.Version,
		barfileOption:        "-b",
		applicationOption:    "-k",
		deployAsSourceOption: "-deployAsSource",
		traceOption:          "-trace",
		verboseOption:        "-v",
		overrideOption:       "-o",
		overridePropOption:   "-p",
		overridesFile:        config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\iibEnvSp\\" + "dev.properties",
		overrideBarFilePath:  config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\target\\iib-overrides\\",
		overrideBarFileName:  config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\target\\iib-overrides\\" + config.Project.ArtifactID + "-" + config.Project.Version + ".bar",
		tracePath:            config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\" + "target\\",
		traceFile:            "createbar.txt",
		cleanBuildOption:     "-cleanBuild",
		barfileName:          config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\" + "target\\" + config.Project.ArtifactID + "-" + config.Project.Version + ".bar",
	}

	fmt.Println(mqsiBarOverride.mqsi,
		mqsiBarOverride.barfileOption,
		mqsiBarOverride.barfileName,
		mqsiBarOverride.applicationOption,
		mqsiBarOverride.artifactID,
		mqsiBarOverride.overrideOption,
		mqsiBarOverride.overrideBarFileName,
		mqsiBarOverride.overridePropOption,
		mqsiBarOverride.overridesFile,
		mqsiBarOverride.verboseOption,
		mqsiBarOverride.tracePath+mqsiBarOverride.traceFile,
	)

	//Create the target directory of iib-overrides if it doesnt exits as the mqsibaroverride command doesnt create it on its own
	CreateDirIfNotExist(mqsiBarOverride.overrideBarFilePath)

	// mqsiapplybaroverride can be applied multiple time based on how many environment specific files
	// are in the config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\iibEnvSp\\" folder

	//Get all the file under the config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\iibEnvSp\\" folder

	files, walkerr := ioutil.ReadDir(config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\iibEnvSp")
	if walkerr != nil {
		log.Fatal(err)
	}

	//Run a loop to all the environment specific files.

	for _, file := range files {
		fmt.Println(file.Name())
		// Reset Override file name to each file
		mqsiBarOverride.overridesFile = config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\iibEnvSp\\" + file.Name()
		mqsiBarOverride.overrideBarFileName = config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\target\\iib-overrides\\" + config.Project.ArtifactID + "-" + config.Project.Version + "-" + strings.Split(file.Name(), ".")[0] + ".bar"

		fmt.Println(mqsiBarOverride.mqsi,
			mqsiBarOverride.barfileOption,
			mqsiBarOverride.barfileName,
			mqsiBarOverride.applicationOption,
			mqsiBarOverride.artifactID,
			mqsiBarOverride.overrideOption,
			mqsiBarOverride.overrideBarFileName,
			mqsiBarOverride.overridePropOption,
			mqsiBarOverride.overridesFile,
			mqsiBarOverride.verboseOption,
			mqsiBarOverride.tracePath+mqsiBarOverride.traceFile,
		)

		cmd := exec.Command(mqsiBarOverride.mqsi,
			mqsiBarOverride.barfileOption,
			mqsiBarOverride.barfileName,
			mqsiBarOverride.applicationOption,
			mqsiBarOverride.artifactID,
			mqsiBarOverride.overrideOption,
			mqsiBarOverride.overrideBarFileName,
			mqsiBarOverride.overridePropOption,
			mqsiBarOverride.overridesFile,
			mqsiBarOverride.verboseOption,
			mqsiBarOverride.tracePath+mqsiBarOverride.traceFile)

		out, err := cmd.Output()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		fmt.Printf("%s", out)
	}

	//start reading the createbar.txt for more verbose output

	fmt.Printf("%s", nil)
	return nil, nil

}

//CreateDirIfNotExist will create the passed dir string
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
