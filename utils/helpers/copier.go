package helpers

import (
	"archive/zip"
	"io/ioutil"
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

func InArray(item string, arr []string) bool {
	for _, check :=  range arr {
		if item == check {
			return true
		}
	}

	return false
}