package settings

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Settings struct {
	Repository Repository `yaml:"repository"`
}

type Repository struct {
	Folder string `yaml:"folder"`
	Branch string `yaml:"branch"`
}

// Open reads already existing configuration from YAML file
func (settings *Settings) Open(filename string) (*Settings, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(filename)
			if err != nil {
				fmt.Println("Couldn't create file")
				panic(err)
			}
		}
		if os.IsPermission(err) {
			fmt.Println("No permission")
			panic(err)
		}
	}

	var config *Settings
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	return config, nil
}

// Save saves user configuration in YAML file
func (settings *Settings) Save(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		panic(err)
	}

	b, err := yaml.Marshal(&settings)
	ioutil.WriteFile(filename, b, 0644)

	return true, nil
}
