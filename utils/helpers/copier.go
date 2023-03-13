package helpers

import (
	"archive/zip"
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadFile(file *zip.File) ([]byte, error) {
	fc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fc.Close()

	content, err := ioutil.ReadAll(fc)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func CreateFile(fileName string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
    if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}