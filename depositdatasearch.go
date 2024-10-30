package main

import (
	"bufio"
	"fmt"
	"github.com/CryptexWebDev/Deposit-Send/common"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func _searchFile(path string) (fileNames []string, err error) {
	fileMask := filepath.Base(path)
	dir := strings.Replace(path, fileMask, "", -1)
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
				//log.Debug("- path:", path, ", name:", d.Name())
				if mask.MatchString(d.Name()) {
					fileNames = append(fileNames, path)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		return []string{path}, nil
	}
	return
}

func SelectDepositData() (index int, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Select deposit data: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	return
}
