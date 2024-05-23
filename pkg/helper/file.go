package helper

import (
	"fmt"
	"os"
	"path/filepath"
)

func FileRead(path string) ([]byte, error) {
	projectDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current working directory. error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(projectDir, path))
	if err != nil {
		return nil, fmt.Errorf("error reading file. error: %v", err)
	}
	return data, nil
}
