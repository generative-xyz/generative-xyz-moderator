package helpers

import (
	"encoding/json"
	"os"
)


func WriteFile(fileName string, data interface{}) {
	f, err := os.Create(fileName)
	if err != nil {
		return 
	}
	byteData,  err := json.Marshal(data)
	if err != nil {
		return 
	}

	_, err = f.Write(byteData)
	if err != nil {
		return 
	}
	f.Sync()
	defer f.Close()
}