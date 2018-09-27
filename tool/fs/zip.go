package fs

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//From https://stackoverflow.com/questions/49057032/recursively-zipping-a-directory-in-golang
func RecursiveZip(pathToZip, destinationPath string) error {
	//destinationFile, err := os.Create(destinationPath)
	destinationFile, err := os.OpenFile(destinationPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}
