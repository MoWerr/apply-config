package data

import (
	"encoding/json"
	"fmt"
	"os"
)

const globalKey = "global"

type Instance map[string]interface{}
type Data map[string]Instance

func Read(dataJSON []byte) (Data, error) {
	data := make(Data)
	err := json.Unmarshal(dataJSON, &data)
	return data, err
}

func ReadFile(filename string) (Data, error) {
	json, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return Read(json)
}

func ReadFiles(filenames ...string) (Data, error) {
	d := make(Data)
	for _, filename := range filenames {
		t, err := ReadFile(filename)
		if err != nil {
			return nil, err
		}
		err = d.merge(t)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (r Data) GetInstanced(instanceName string) Instance {
	l := len(r[globalKey]) + len(r[instanceName])
	d := make(Instance, l)
	d.merge(r[globalKey])
	d.merge(r[instanceName])
	return d
}

func (r Data) merge(other Data) error {
	for k, v := range other {
		if _, ok := r[k]; !ok {
			r[k] = make(Instance, len(other[k]))
		}
		err := r[k].merge(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r Instance) merge(other Instance) error {
	for k, v := range other {
		if _, ok := r[k]; ok {
			return fmt.Errorf("%q data key is duplicated", k)
		}
		r[k] = v
	}
	return nil
}
