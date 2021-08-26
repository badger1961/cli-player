package common

import (
	"errors"
	"os"
)

func CheckInputFile(fileName string) error {
	fileInfo, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		return errors.New("Hmm ... File " + fileName + " not found")
	}

	if fileInfo.IsDir() {
		return errors.New("Hmm ... File " + fileName + " should be a file not folder")
	}

	return nil
}
