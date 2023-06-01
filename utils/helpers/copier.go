package helpers

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
)

func Openfile(filename string) ([]byte, error) {
	fc, err := os.Open(filename)
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

func CreateTxtFile(fileName string, bytes []byte) error {
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

func ZipFiles(fileName string, zipFileNames ...string) (*os.File, error) {
	fmt.Println("creating zip archive...")
	archive, err := os.Create(fmt.Sprintf("%s.zip", fileName))
	if err != nil {
		return nil, err
	}

	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	for _, zipFile := range zipFileNames {
		f1, err := os.Open(zipFile)
		if err != nil {
			return archive, nil
		}
		defer f1.Close()

		fmt.Println("writing first file to archive...")
		w1, err := zipWriter.Create(fmt.Sprintf("%s/%s", fileName, zipFile))
		if err != nil {
			return archive, nil
		}

		if _, err := io.Copy(w1, f1); err != nil {
			return archive, nil
		}
	}

	return archive, nil

}

func FileExists(filename string) (*fs.FileInfo, error) {
	fInfor, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	return &fInfor, nil
}
