package command

import (
	"errors"
	"log"
	"os"
	"path"
	"strings"
)

func PathIsValid(path string) error {
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

func SplitPath(cpath string) (string, string, string) {
	ext := path.Ext(cpath)
	splits := strings.Split(cpath, "/")
	filename := splits[len(splits)-1]
	cpath = strings.TrimSuffix(cpath, filename)
	filename = strings.TrimSuffix(filename, ext)
	return cpath, filename, ext
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func InitEnvs() {
	GlobalOption.EcoURL = os.Getenv("ECO_URL")
	if GlobalOption.EcoURL == "" {
		log.Fatal("env ECO_URL is required")
	}
	GlobalOption.APIKey = os.Getenv("ECO_TOKEN")
	if GlobalOption.APIKey == "" {
		log.Fatal("env ECO_TOKEN is required")
	}
	GlobalOption.ProjectID = os.Getenv("ECO_PROJECT")
	if GlobalOption.ProjectID == "" {
		log.Fatal("env ECO_PROJECT is required")
	}
	GlobalOption.RegionID = os.Getenv("ECO_REGION")
	if GlobalOption.RegionID == "" {
		log.Fatal("env ECO_REGION is required")
	}
	GlobalOption.UserID = os.Getenv("ECO_USER_ID")
	if GlobalOption.UserID == "" {
		log.Fatal("env ECO_USER_ID is required")
	}

}
