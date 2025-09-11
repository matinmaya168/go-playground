package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProtoMergeOptions struct {
	FolderPath       string
	OutPutFilePath   string
	IgnoreFiles      []string
	PkgFinder        string
	PkgReplacer      string
	OptGoPkgFinder   string
	OptGoPkgReplacer string
}

func MergeProto(pmo ProtoMergeOptions) {
	files, err := os.ReadDir(pmo.FolderPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var fileContents []byte
	var startContents []byte
	for _, file := range files {
		if file.IsDir() || nameExists(file.Name(), pmo.IgnoreFiles) {
			continue
		}

		fileName := filepath.Join(pmo.FolderPath, file.Name())
		fileData, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Printf("Error reading file %s: %s\n", fileName, err)
			continue
		}

		startContents = fileData[:69]
		fileData = fileData[70:]
		fileContents = append(fileContents, fileData...)
	}

	startContents = append(startContents, fileContents...)

	start := strings.Replace(string(startContents), pmo.PkgFinder, pmo.PkgReplacer, 1)
	start = strings.Replace(start, pmo.OptGoPkgFinder, pmo.OptGoPkgReplacer, 1)
	err = os.WriteFile(pmo.OutPutFilePath, []byte(start), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Merged files to:", pmo.OutPutFilePath)
}

func nameExists(name string, names []string) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Merge proto files running...")
}
