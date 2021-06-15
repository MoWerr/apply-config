package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/MoWerr/apply-config/config"
	"github.com/MoWerr/apply-config/data"
	"github.com/MoWerr/apply-config/templates"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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

	templatePaths := config.GetTemplatePaths()
	templates, err := templates.ReadFiles(templatePaths...)
	check(err)

	err = templates.Deploy(*config, data)
	check(err)
}
