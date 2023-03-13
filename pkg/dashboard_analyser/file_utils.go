package dashboard_analyser

import (
	"os"
)

// read the file  with the given filename in buffer
func loadFile(filename string) ([]byte, error) {

	buffer, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
