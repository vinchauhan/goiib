package mqsi

type Build struct {
	Project struct {
		ModelVersion string `yaml:"modelVersion"`
		GroupID      string `yaml:"groupId"`
		ArtifactID   string `yaml:"artifactId"`
		Version      string `yaml:"version"`
		Packaging    string `yaml:"packaging"`
		Profiles     struct {
			Profile struct {
				ID         string `yaml:"id"`
				Activation struct {
					ActiveByDefault bool `yaml:"activeByDefault"`
				} `yaml:"activation"`
				Properties struct {
					Workspace                          string `yaml:"workspace"`
					InitialDeletes                     string `yaml:"initialDeletes"`
					UnpackIibDependenciesIntoWorkspace bool   `yaml:"unpackIibDependenciesIntoWorkspace"`
					FailOnInvalidProperties            bool   `yaml:"failOnInvalidProperties"`
					UseClassloaders                    bool   `yaml:"useClassloaders"`
					FailOnInvalidClassloader           bool   `yaml:"failOnInvalidClassloader"`
					CreateOrPackageBar                 string `yaml:"createOrPackageBar"`
					CompleteDeployment                 bool   `yaml:"completeDeployment"`
					TimeoutSecs                        int    `yaml:"timeoutSecs"`
					MqsiCreateBarDeployAsSource        bool   `yaml:"mqsiCreateBarDeployAsSource"`
				} `yaml:"properties"`
			} `yaml:"profile"`
			Dependencies struct {
				GroupID    string `yaml:"groupId"`
				ArtifactID string `yaml:"artifactId"`
				Version    string `yaml:"version"`
			} `yaml:"dependencies"`
		} `yaml:"profiles"`
	} `yaml:"project"`
}
