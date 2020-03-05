package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	print(out, path, true, false, "", printFiles)
	return nil
}

func print(out io.Writer, path string, doNotPrint bool, isElemLast bool, indent string, printFiles bool) {
	file, _ := os.Open(path)
	defer file.Close()

	fileInfo, _ := file.Stat()

	if !doNotPrint {
		sym, _ := lineSyms[isElemLast]
		fmt.Fprint(out, indent + sym, fileInfo.Name(), fileSize(fileInfo), "\n")
		delimiter, _ := tabSyms[isElemLast]
		indent = indent + delimiter
	}
	if !fileInfo.IsDir() {
		return
	}

	fileInfos, _ := file.Readdir(0)
	if !printFiles {
		fileInfos = filterDirs(fileInfos)
	}
	sort.SliceStable(fileInfos, func(i, j int) bool { return fileInfos[i].Name() < fileInfos[j].Name() })
	for i, fi := range fileInfos {
		isLast := i == len(fileInfos) - 1

		print(out, filepath.Join(path, fi.Name()), false, isLast, indent, printFiles)
	}
}

var lineSyms = map[bool]string{
	false: "├───",
	true:  "└───",
}
var tabSyms = map[bool]string{
	false: "│\t",
	true:  "\t",
}

func filterDirs(s []os.FileInfo) []os.FileInfo {
	var tmp []os.FileInfo
	for _, elem := range s {
		if elem.IsDir() {
			tmp = append(tmp, elem)
		}
	}
	return tmp
}

func fileSize(fi os.FileInfo) string {
	if fi.IsDir() {
		return ""
	}
	if fi.Size() > 0 {
		return fmt.Sprintf(" (%db)", fi.Size())
	}
	return " (empty)"
}
