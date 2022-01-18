package zfile

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// ReadAll read file all content.
func ReadAll(path string) ([]byte, error) {
	fd, err := os.Open(path)
	defer Close(fd)
	if err != nil {
		return nil, err
	}

	fileData, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

// ReadBlock read file with line.
func ReadBlock(fd *os.File, offset int64, size int) ([]byte, int, bool, error) {
	if fd == nil {
		return nil, 0, false, fmt.Errorf("fd is nil")
	}

	_, err := fd.Seek(offset, 0)
	if err != nil {
		return nil, 0, false, err
	}
	buffer := make([]byte, size)

	readSize, err := fd.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return buffer, readSize, true, nil
		}

		return nil, readSize, false, err
	}

	return buffer, readSize, false, nil
}

// ReadLine read file with line, index begin from 0.
func ReadLine(fd *os.File, index int64) ([]byte, error) {
	if fd == nil {
		return nil, fmt.Errorf("fd is nil")
	}

	if index < 0 {
		return make([]byte, 0, 0), nil
	}

	var curIndex int64
	fileScanner := bufio.NewScanner(fd)
	for fileScanner.Scan() {
		if curIndex == index {
			return fileScanner.Bytes(), nil
		}

		curIndex++
	}

	return make([]byte, 0, 0), nil
}
