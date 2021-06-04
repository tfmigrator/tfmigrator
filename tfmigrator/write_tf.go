package tfmigrator

import (
	"fmt"
	"io"
	"io/ioutil"
)

// WriteTFInTemporalFile writes Terraform Configuration in a temporal file.
// In this function, a temporal file is created and the path to the file is returned, so you have to remove the file by yourself.
func WriteTFInTemporalFile(reader io.Reader) (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf("create a temporal file to write Terraform configuration (.tf): %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return f.Name(), fmt.Errorf("write standard input to a temporal file %s: %w", f.Name(), err)
	}
	return f.Name(), nil
}
