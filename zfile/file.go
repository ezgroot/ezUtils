package zfile

import (
	"fmt"
	"io/ioutil"
	"os"
)

// IsExist check is file or path exist, and return file info.
func IsExist(path string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		return fileInfo, nil
	}

	if os.IsNotExist(err) {
		return fileInfo, os.ErrNotExist
	}

	return fileInfo, err
}

// GetAllFile get all file list in this path.
func GetAllFile(path string) error {
	rd, err := ioutil.ReadDir(path)
	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("[%s]\n", path+"/"+fi.Name())
			GetAllFile(path + fi.Name() + "\\")
		} else {
			fmt.Println(path + "/" + fi.Name())
		}
	}

	return err
}

// Close close file handle.
func Close(file *os.File) {
	if file == nil {
		return
	}

	err := file.Close()
	if err != nil {
		fmt.Printf("file close error = %s", err)
	}
}
