package main

import (
	"bufio"
	"fmt"
	"github.com/CryptexWebDev/Deposit-Send/common"
	"os"
	"path/filepath"
	"strings"
)

/*
`^deposit_data-\d{10}\.json$`gm
*/
func _searchFile(path string) (fileNames []string, err error) {
	fileMask := filepath.Base(path)
	dir := strings.Replace(path, fileMask, "", -1)
	if !common.IsFileExists(dir) {
		return nil, fmt.Errorf("directory not exists: %s", dir)
	}
	if strings.Index(fileMask, "*") != 0 {
		fileMask = "^" + strings.Replace(fileMask, "*", "\\d{10}\\", 1) + "$"
		//mask, err := regexp.Compile(fileMask)
		if err != nil {
			return nil, err
		}
	}
	println("dir,fileMask:", dir, fileMask)
	if dir == "" {
		dir = "."
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
