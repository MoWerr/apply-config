package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

type FileConfig struct {
	Source      string
	Destination string
}

type InstanceConfig struct {
	Name        string
	Destination string
}

type Config struct {
	Files     []FileConfig
	Instances []InstanceConfig
}

func main() {
	config := Config{}

	configJson, _ := ioutil.ReadFile("apply-config.json")
	json.Unmarshal(configJson, &config)

	templateValues := make(map[string]map[string]interface{})

	files, _ := ioutil.ReadDir(".")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") && !strings.EqualFold(file.Name(), "apply-config.json") {
			valuesJson, err := ioutil.ReadFile(file.Name())
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(valuesJson, &templateValues)
			if err != nil {
				panic(err)
			}
		}
	}

	for _, file := range config.Files {
		t, _ := template.ParseFiles(file.Source)

		for _, instance := range config.Instances {
			fileDest := path.Join(instance.Destination, file.Destination)
			fileDirPath := path.Dir(fileDest)

			instanceData := make(map[string]interface{}, len(templateValues["global"])+len(templateValues[instance.Name]))
			for k, v := range templateValues["global"] {
				instanceData[k] = v
			}
			for k, v := range templateValues[instance.Name] {
				instanceData[k] = v
			}

			err := os.MkdirAll(fileDirPath, 0744)
			if err != nil {
				fmt.Println(err)
				return
			}

			fileWriter, err := os.Create(fileDest)
			if err != nil {
				fmt.Println(err)
				return
			}

			defer fileWriter.Close()

			err = t.Execute(fileWriter, instanceData)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
