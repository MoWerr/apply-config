package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
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
	configPath := flag.String("c", "config.json", "Specifies path to the main configuration file.")
	dataPaths := flag.String("d", "data.json", "Colon separated paths to the data files. It will be considered only if there are no data sources defined in the config file.")
	testRun := flag.Bool("t", false, "Allows to perform test run. It (will print generated files in standard output).")
	outputPath := flag.String("o", "", "Can be specified only along with 't'. It will deploy test files to specified location instead of printing them to the standard output.")

	flag.Parse()

	if *outputPath != "" && !*testRun {
		check(errors.New("-o flag requires also -t to work"))
		fmt.Println("")
		flag.PrintDefaults()
		os.Exit(1)
	}

	config, err := config.ReadFile(*configPath, strings.Split(*dataPaths, ":"))
	check(err)

	data, err := data.ReadFiles(config.DataSources...)
	check(err)

	templatePaths := config.GetTemplatePaths()
	templates, err := templates.ReadFiles(templatePaths...)
	check(err)

	if *testRun && *outputPath == "" {
		err = templates.Deploy(*config, data, deployStdout, nil)
	} else {
		err = templates.Deploy(*config, data, deployFile, *outputPath)
	}

	check(err)
}

func deployStdout(deployer templates.Deployer, instance config.Instance, file config.File, _ interface{}) error {
	dest := path.Join(instance.Destination, file.Destination)

	fmt.Println("======================================================================")
	fmt.Printf("## Instance %q \n", instance.Name)
	fmt.Printf("## File %q \n", dest)
	fmt.Println("======================================================================")

	err := deployer.Deploy(os.Stdout)
	if err != nil {
		return fmt.Errorf("Failed to print the file: %q; %w", dest, err)
	}

	fmt.Println("\n======================================================================")
	fmt.Println()
	fmt.Println()

	return nil
}

func deployFile(deployer templates.Deployer, instance config.Instance, file config.File, customLocation interface{}) error {
	location := customLocation.(string)
	dest := ""

	if location == "" {
		dest = path.Join(instance.Destination, file.Destination)
	} else {
		dest = path.Join(location, instance.Name, file.Destination)
	}

	destPath := path.Dir(dest)

	err := os.MkdirAll(destPath, 0766)
	if err != nil {
		return fmt.Errorf("Failed to make destination path: %q; %w", destPath, err)
	}

	w, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("Failed to create the destination file: %q; %w", dest, err)
	}

	defer w.Close()

	err = deployer.Deploy(w)
	if err != nil {
		return fmt.Errorf("Failed to deploy the file: %q; %w", dest, err)
	}

	return nil
}
