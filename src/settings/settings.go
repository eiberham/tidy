package settings

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Settings struct {
	Git struct {
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"database"`
	Repository struct {
		Folder string `yaml:"folder"`
		Branch string `yaml:"branch"`
	} `yaml:"database"`
}

func (settings *Settings) Open(filename string) (*Settings, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create("/tmp/settings.yaml")
		if err != nil {
			panic(err)
		}
		/* config Settings
		return config{
			Git: {
				User: ""
				Pass: ""
			}
			Repository: {
				Folder: ""
				Branch: ""
			}
		} */
	}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//var settings Settings
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&settings)
	if err != nil {
		panic(err)
	}

	return settings, nil
}

func (settings *Settings) Save(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		panic(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(&settings)
	if err != nil {
		return false, err
	}

	return true, nil
}
