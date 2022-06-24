package utils

import (
	"io/fs"
	"os"
	"path/filepath"
)

func GetAbsPath() string {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return path
}

func GetExtFiles(path string, ext string) []string {
	/**
	获取路径下包含子路径下所有的文件
	*/
	var files []string
	err := filepath.Walk(path, func(file string, info fs.FileInfo, err error) error {
		if !info.IsDir() && ext == filepath.Ext(file) {
			files = append(files, file)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func ScanDir(path string) []string {
	/**
	获取路径下包含子路径下所有的文件
	*/
	var files []string
	err := filepath.Walk(path, func(file string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, file)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
