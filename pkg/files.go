package jt

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

// read json contents of a file into target struct
func ReadFile(path string, target interface{}) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Errorf("file read error, %v", err)
	}

	err = json.Unmarshal(f, &target)
	if err != nil {
		return errors.Errorf("json unmarshal error, %v", err)
	}
	return nil
}

// write source struct as json to stdout or file, if provided
func WriteFile(path string, source interface{}) error {
	bites, err := json.MarshalIndent(source, "", "  ")
	if err != nil {
		return errors.Errorf("json marshal error, %v", err)
	}

	bites = append(bites, byte('\n'))

	if path != "" {
		err = ioutil.WriteFile(path, bites, 0644)
		if err != nil {
			errors.Errorf("file write error, %v", err)
		}
	} else {
		_, err := os.Stdout.Write(bites)
		if err != nil {
			errors.Errorf("stdout write error, %v", err)
		}
	}

	return nil
}
