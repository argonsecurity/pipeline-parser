package testutils

import (
	"io/ioutil"
)

func ReadFile(filepath string) []byte {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	return data
}
