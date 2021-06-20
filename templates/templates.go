package templates

import (
	"fmt"
	"io"
	"os"
	"path"
	"text/template"

	"github.com/MoWerr/apply-config/config"
	"github.com/MoWerr/apply-config/data"
)

type Deployer struct {
	template *template.Template
	data     data.Instance
}

func (d *Deployer) Deploy(writer io.Writer) error {
	return d.template.Execute(writer, d.data)
}

type Templates template.Template
type DeployFunc func(Deployer, config.Instance, config.File, interface{}) error

func ReadFiles(filenames ...string) (*Templates, error) {
	t := template.New("")

	for _, filename := range filenames {
		f, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("Failed to read the template file: %q; %w", filename, err)
		}

		_, err = t.New(filename).Parse(string(f))
		if err != nil {
			return nil, fmt.Errorf("Failed to parse the template file: %q; %w", filename, err)
		}
	}

	return (*Templates)(t), nil
}

func (r *Templates) Deploy(config config.Config, data data.Data, deployFunc DeployFunc, userData interface{}) error {
	return r.DeployInstances(nil, config, data, deployFunc, userData)
}

func (r *Templates) DeployInstances(instances []string, config config.Config, data data.Data, deployFunc DeployFunc, userData interface{}) error {
	for _, instance := range config.Instances {
		if instances == nil || containsInstance(instances, instance.Name) {
			err := r.deployInstance(instance, config.Files, data, deployFunc, userData)
			if err != nil {
				return fmt.Errorf("Failed to deploy the instance: %q; %w", instance.Name, err)
			}
		}
	}

	return nil
}

func containsInstance(instances []string, instance string) bool {
	for _, i := range instances {
		if i == instance {
			return true
		}
	}

	return false
}

func (r *Templates) deployInstance(instance config.Instance, files []config.File, data data.Data, deployFunc DeployFunc, userData interface{}) error {
	for _, file := range files {
		t := (*template.Template)(r)
		deployer := Deployer{
			template: t.Lookup(file.Source),
			data:     data.GetInstanced(instance.Name),
		}

		dest := path.Join(instance.Destination, file.Destination)
		err := deployFunc(deployer, instance, file, userData)
		if err != nil {
			return fmt.Errorf("Failed to deploy the file: %q; %w", dest, err)
		}
	}

	return nil
}
