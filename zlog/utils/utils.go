package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// OS types.
const (
	Linux   = "linux"
	Windows = "windows"
	Darwin  = "darwin"
)

// GetOS get os type
func GetOS() (os string) {
	return runtime.GOOS
}

// IsExist is file or dir exists
func IsExist(path string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		return fileInfo, nil
	}

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("%s not exist", path)
	}

	return nil, err
}

// GetJSONStr get json string
func GetJSONStr(key string, value interface{}) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var infoMap = make(map[string]interface{})
	infoMap[key] = value

	b, err := json.Marshal(&infoMap)
	if err != nil {
		fmt.Printf("err = %s\n", err)
	}

	t1 := strings.Trim(string(b), "{")

	return strings.Trim(t1, "}")
}
