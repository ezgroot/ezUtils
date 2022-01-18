package zfile

import (
	"io/ioutil"
	"os"
)

// WriteCoverFromHead From the beginning of the file, write over.
func WriteCoverFromHead(path string, data []byte) error {
	fileHandle, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	defer Close(fileHandle)
	if err != nil {
		return err
	}

	_, err = fileHandle.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = fileHandle.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// WriteCoverAll Delete all the contents of the file then write from the beginning.
func WriteCoverAll(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// WriteAppendToTail Appends writing to the end of the file.
func WriteAppendToTail(path string, data []byte) error {
	fileHandle, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	defer Close(fileHandle)
	if err != nil {
		return err
	}

	_, err = fileHandle.Seek(0, 2)
	if err != nil {
		return err
	}

	_, err = fileHandle.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// WriteNewLineAppendToTail Appends writing to the end of the file.
func WriteNewLineAppendToTail(path string, data []byte) error {
	fileHandle, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	defer Close(fileHandle)
	if err != nil {
		return err
	}

	_, err = fileHandle.Seek(0, 2)
	if err != nil {
		return err
	}

	_, err = fileHandle.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	_, err = fileHandle.Write(data)
	if err != nil {
		return err
	}

	return nil
}
