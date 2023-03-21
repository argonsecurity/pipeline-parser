package utils

import "os"

func ReadFile(path string) ([]byte, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
