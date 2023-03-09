package helpers

import (
	"os"
)

func ReadFile(fileName string) ([]byte, error) {
	f, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func CreateFile(fileName string, data []byte) error {
	f, err := os.Create(fileName)
    if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}