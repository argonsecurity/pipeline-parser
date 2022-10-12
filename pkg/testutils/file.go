package testutils

import (
	"os"
)

func ReadFile(filepath string) []byte {
	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	return data
}
