package templates

import (
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
			return nil, err
		}
		_, err = t.New(filename).Parse(string(f))
		if err != nil {
			return nil, err
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
		return err
	}
	w, err := os.Create(f)
	if err != nil {
		return err
	}
	defer w.Close()
	t := (*template.Template)(r)
	return t.ExecuteTemplate(w, file.Source, data.GetInstanced(instance.Name))
}
