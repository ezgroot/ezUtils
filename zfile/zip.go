package zfile

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Zip(dstParent string, dstFile string, srcParent string, srcFile string) error {
	dstFileFullPath := fmt.Sprintf("%s/%s", strings.TrimSuffix(dstParent, "/"), dstFile)
	srcFileFullPath := fmt.Sprintf("%s/%s", strings.TrimSuffix(srcParent, "/"), srcFile)

	fw, err := os.Create(dstFileFullPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	zw := zip.NewWriter(fw)
	defer func() {
		if err := zw.Close(); err != nil {
			fmt.Printf("zip Close error = %s\n", err)
		}
	}()

	return filepath.Walk(srcFileFullPath, func(path string, fi os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}

		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		trimPath := strings.TrimPrefix(path, srcParent)

		fh.Name = strings.TrimPrefix(trimPath, string(filepath.Separator))

		// if not has this, while unzip it will not be a dir
		if fi.IsDir() {
			fh.Name += "/"
		}

		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		// check if it is a dir or a regular file, we only write file head info, not write content
		if !fh.Mode().IsRegular() {
			return nil
		}

		fr, err := os.Open(path)
		if err != nil {
			return err
		}

		n, err := io.Copy(w, fr)
		if err != nil {
			fr.Close()
			return err
		}

		fr.Close()

		fmt.Printf("success zip： %s, total write %d number byte content\n", path, n)

		return nil
	})
}

func UnZip(srcFileFullPath string, dstParent string) (err error) {
	zr, err := zip.OpenReader(srcFileFullPath)
	if err != nil {
		return
	}
	defer zr.Close()

	if dstParent != "" {
		if err := os.MkdirAll(dstParent, 0755); err != nil {
			return err
		}
	}

	for _, file := range zr.File {
		path := filepath.Join(dstParent, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}

			continue
		}

		fr, err := file.Open()
		if err != nil {
			return err
		}

		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			fr.Close()
			return err
		}

		n, err := io.Copy(fw, fr)
		if err != nil {
			fw.Close()
			fr.Close()
			return err
		}

		fw.Close()
		fr.Close()

		fmt.Printf("success unzip： %s, total write %d number byte content\n", path, n)
	}

	return nil
}

func ParseZipInMemory(parent string, file string) (*zip.Reader, error) {
	fileFullPath := fmt.Sprintf("%s/%s", strings.TrimSuffix(parent, "/"), file)

	fileInfo, err := os.Stat(fileFullPath)
	if err != nil {
		return nil, err
	}

	fd, err := os.Open(fileFullPath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	zipReader, err := zip.NewReader(fd, fileInfo.Size())
	if err != nil {
		return nil, err
	}

	for i, zipMember := range zipReader.File {
		if zipMember.FileInfo().IsDir() {
			fmt.Printf("formfile[%d]: dirname=[%s]\n", i, zipMember.Name)
		} else {
			fmt.Printf("formfile[%d]: filename=[%s]\n", i, zipMember.Name)

			f, err := zipMember.Open()
			if err != nil {
				return nil, err
			}

			buf, err := ioutil.ReadAll(f)
			if err != nil {
				f.Close()
				fmt.Printf("ioutil.ReadAll error = %s\n", err)
				return nil, err
			}

			f.Close()

			fmt.Printf("formfile[%d]: filename=[%s], size=%d, content=[%s]\n", i, zipMember.Name, len(buf), strings.TrimSuffix(string(buf), "\n"))
		}
	}

	return zipReader, nil
}
