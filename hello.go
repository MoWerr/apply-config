package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/MoWerr/apply-config/config"
	"github.com/MoWerr/apply-config/data"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	configPath := flag.String("c", "config.json", "Specifies path to the main configuration file")
	dataPaths := flag.String("d", "data.json", "Colon separated paths to the data files")

	flag.Parse()

	config, err := config.ReadFile(*configPath, strings.Split(*dataPaths, ":"))
	check(err)

	data, err := data.ReadFiles(config.DataSources...)
	check(err)

	fmt.Println(data.GetInstanced("artemis"))
	fmt.Println(data.GetInstanced("beta"))
	/*

		files, _ := ioutil.ReadDir(".")
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".json") && !strings.EqualFold(file.Name(), "app-config.json") {
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
	*/
}
