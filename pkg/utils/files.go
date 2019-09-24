package utils

import "io/ioutil"

//WriteFile writes to a file
func WriteFile(filepath string, data []byte) error {
	return ioutil.WriteFile(filepath, data, 0644)
}
