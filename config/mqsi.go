package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

//CompileProject -cleanBuild the IIB project and creates a BAR file in the /target directory
func CompileProject() ([]byte, error) {

	log.Infof("Scanning for projects...")
	log.Infof("")
	config := BuildConfig{}
	path := filepath.Join(filepath.Dir("."), "build.yaml")
	//fmt.Println(path)

	source, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to %v", err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return nil, fmt.Errorf("Could not Unmarshal the build %v", err)
	}

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

	log.Infof("------------------------------------------------------------------------")
	log.Infof("Building " + mqsiCreateBar.artifactID + " " + mqsiCreateBar.version)
	log.Infof("------------------------------------------------------------------------")
	log.Infof("")
	log.Infof("")

	log.Infof("Checking for dependencies... ")
	log.Infof("")

	//Dependencies implementations

	log.Infof("Looking for .project file %s", filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, ".project"))
	log.Infof("")

	//projectFileBytes, err := ioutil.ReadFile(filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, ".project"))

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
		filepath.Join(mqsiCreateBar.tracePath, "createbartrace.txt"),
	)

	// Set default shell and option for linux
	shell, option, cmdSeperator := "/bin/bash", "-c", ";"
	//Detect if Windows
	if runtime.GOOS == "windows" {

		log.Infof("Detected Windows environment")
		shell, option, cmdSeperator = "cmd", "/c", "&&"

	}

	finalCmd := fmt.Sprintln("mqsiprofile", cmdSeperator, mqsiCreateBarCmdString)

	fmt.Println("")
	fmt.Println("generated mqsiCommand follows...")
	fmt.Println("---------------------")

	fmt.Printf("%s", finalCmd)

	fmt.Println("---------------------")

	cmd := exec.Command(shell, option, finalCmd)
	cmd.Stderr = os.Stderr

	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		fmt.Errorf("Couldnt not execute mqsicreatebar : %v", err)
	}

	log.Infof("-----------------------------------------------")
	log.Infof("Creating overrides properties for the developer")
	log.Infof("-----------------------------------------------")
	fmt.Println("")

	//Create defalt.properties file in the /target directory to use for creating the override file.
	mqsiReadBarCmd := fmt.Sprintln("mqsireadbar",
		"-b",
		mqsiCreateBar.barfileName,
		"-r",
	)

	fmt.Printf("mqsiReadBar command to create default.properties for overide is : %s", mqsiReadBarCmd)

	fmt.Println("")
	fmt.Println("")

	//start reading the createbar.txt for more verbose output
	finalmqsiReadBarCmd := fmt.Sprintln("mqsiprofile", cmdSeperator, mqsiReadBarCmd)

	defaultPropCmd := exec.Command(shell, option, finalmqsiReadBarCmd)

	defaultout, defaulterr := defaultPropCmd.Output()
	if defaulterr != nil {
		return nil, fmt.Errorf("Could execute the command mqsireadbar and failed with %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(mqsiCreateBar.tracePath, "default.properties"), []byte(defaultout), 777)

	if err != nil {
		log.Fatalf("mqsireadbar filed with %s\n", err)
	}

	log.Infof("-----------------------------------------------")
	log.Infof("BUILD SUCCESS")
	log.Infof("-----------------------------------------------")

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
		return nil, fmt.Errorf("Failed to %v", err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return nil, fmt.Errorf("Could not Unmarshal the build %v", err)
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
		overridesFile:        filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, "iibEnvSp", "dev.properties"),
		overrideBarFilePath:  filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, "target", "iib-overrides"),
		overrideBarFileName:  filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, "target", "iib-overrides", config.Project.ArtifactID+"-"+config.Project.Version+".bar"),
		tracePath:            filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, "target"),
		traceFile:            "createbar.txt",
		cleanBuildOption:     "-cleanBuild",
		barfileName:          filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, "target", config.Project.ArtifactID+"-"+config.Project.Version+".bar"),
	}

	//Create the target directory of iib-overrides if it doesnt exits as the mqsibaroverride command doesnt create it on its own
	CreateDirIfNotExist(mqsiBarOverride.overrideBarFilePath)

	// mqsiapplybaroverride can be applied multiple time based on how many environment specific files
	// are in the config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\iibEnvSp\\" folder

	//Get all the file under the config.Project.Profiles.Profile.Properties.Workspace + config.Project.ArtifactID + "\\iibEnvSp\\" folder

	files, walkerr := ioutil.ReadDir(filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, "iibEnvSp"))
	if walkerr != nil {
		fmt.Errorf("Could not walk the path of the iibEnvSp : %v", err)
	}

	//Run a loop to all the environment specific files.

	for _, file := range files {
		fmt.Println(file.Name())
		//Reset Override file name to each file

		mqsiBarOverride.overridesFile = filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, "iibEnvSp", file.Name())
		mqsiBarOverride.overrideBarFileName = filepath.Join(config.Project.Profiles.Profile.Properties.Workspace, config.Project.ArtifactID, "target", "iib-overrides", config.Project.ArtifactID+"-"+config.Project.Version+"-"+strings.Split(file.Name(), ".")[0]+".bar")

		mqsiApplyOverrideCmd := fmt.Sprintln(
			mqsiBarOverride.mqsi,
			mqsiBarOverride.barfileOption,
			mqsiBarOverride.barfileName,
			mqsiBarOverride.applicationOption,
			mqsiBarOverride.artifactID,
			mqsiBarOverride.overrideOption,
			mqsiBarOverride.overrideBarFileName,
			mqsiBarOverride.overridePropOption,
			mqsiBarOverride.overridesFile,
			mqsiBarOverride.verboseOption,
			filepath.Join(mqsiBarOverride.tracePath, mqsiBarOverride.traceFile),
		)

		fmt.Printf("mqsiapplyoverride command is : %s", mqsiApplyOverrideCmd)

		cmd := exec.Command("/bin/bash", "-c", mqsiApplyOverrideCmd)

		out, err := cmd.Output()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %v\n", err)
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
