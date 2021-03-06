package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateConfig(t *testing.T) {

	t.Parallel()

	Convey("Given a build file", t, func() {
		buildFilePath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "vinchauhan", "goiib", "test", "build.yaml")
		Convey("When config object is created", func() {

			config, err := createConfig(buildFilePath)

			if err != nil {
				fmt.Println(err)
				So(err, ShouldBeError)
			}
			//
			Convey("ArtifactID should be CustomerDatabaseV1", func() {
				So(config.Project.ArtifactID, ShouldEqual, "CustomerDatabaseV1")
			})

			Convey("GroupId should be com.springular.iib", func() {
				So(config.Project.GroupID, ShouldEqual, "com.springular.iib")
			})

		})

	})

}
