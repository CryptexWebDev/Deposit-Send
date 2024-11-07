package tools

import (
	"fmt"
	"github.com/CryptexWebDev/Deposit-Send/common"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func IsFileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func FilesSearchByMask(dir, fileMask string) ([]string, error) {
	var fileNames []string
	if !common.IsFileExists(dir) {
		return nil, fmt.Errorf("directory not exists: %s", dir)
	}
	if strings.Index(fileMask, "*") != 0 {
		fileMask = "^" + strings.Replace(fileMask, "*", "\\d{10}\\", 1) + "$"
		mask, err := regexp.Compile(fileMask)
		if err != nil {
			return nil, err
		}
		if dir == "" {
			dir = "."
		}
		err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if path != dir {
				if mask.MatchString(d.Name()) {
					fileNames = append(fileNames, d.Name())
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		return []string{fileMask}, nil
	}
	return fileNames, nil
}
