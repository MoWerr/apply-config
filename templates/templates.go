package templates

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/MoWerr/apply-config/config"
	"github.com/MoWerr/apply-config/data"
)

type Templates template.Template

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

func (r *Templates) Deploy(config config.Config, data data.Data) error {
	for _, instance := range config.Instances {
		for _, file := range config.Files {
			err := r.deployFile(instance, file, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Templates) deployFile(instance config.Instance, file config.File, data data.Data) error {
	f := path.Join(instance.Destination, file.Destination)

	err := os.MkdirAll(path.Dir(f), 0766)
	if err != nil {
		return fmt.Errorf("Failed to make destination path: %q; %w", path.Dir(f), err)
	}

	w, err := os.Create(f)
	if err != nil {
		return fmt.Errorf("Failed to create the destination file: %q; %w", f, err)
	}

	defer w.Close()
	t := (*template.Template)(r)

	err = t.ExecuteTemplate(w, file.Source, data.GetInstanced(instance.Name))
	if err != nil {
		return fmt.Errorf("Failed to apply the template: %q; %w", f, err)
	}

	return nil
}
