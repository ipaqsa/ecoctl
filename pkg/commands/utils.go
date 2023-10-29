package commands

import (
	"errors"
	"os"
	"path"
	"strings"
)

func pathToTemplateIsValid(path string) error {
	splits := strings.Split(path, ".")
	if len(splits) < 2 {
		return errors.New("path invalid")
	}
	format := splits[len(splits)-1]
	if format != "yaml" {
		return errors.New("format is not yaml")
	}
	return nil
}

func splitPath(cpath string) (string, string, string) {
	ext := path.Ext(cpath)
	splits := strings.Split(cpath, "/")
	filename := splits[len(splits)-1]
	cpath = strings.TrimSuffix(cpath, filename)
	filename = strings.TrimSuffix(filename, ext)
	return cpath, filename, ext
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
