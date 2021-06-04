package tfmigrator

import (
	"fmt"
	"io"
	"io/ioutil"
)

func (ctrl *Controller) writeTF() (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf("create a temporal file to write Terraform configuration (.tf): %w", err)
	}
	defer f.Close()
	// read tf from stdin and write a temporal file
	if _, err := io.Copy(f, ctrl.Stdin); err != nil {
		return f.Name(), fmt.Errorf("write standard input to a temporal file %s: %w", f.Name(), err)
	}
	return f.Name(), nil
}
