package utils

import (
	"io"
	"os"
)

// SaveToFile saves the contents from a reader to a file specified by filePath.
//
// Parameters:
// - reader: the input reader containing the data to be saved.
// - filePath: the path to the file where the data will be saved.
// Return type: error, returns any error that occurred during the saving process.
func SaveToFile(reader io.Reader, filePath string) error {
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func(outFile *os.File) {
		_ = outFile.Close()
	}(outFile)

	_, err = io.Copy(outFile, reader)
	if err != nil {
		return err
	}

	return nil
}
