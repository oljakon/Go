package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func printTree(isLast bool, prefix string, info os.FileInfo) (lines string) {
	line := prefix
	if isLast {
		line += "└───"
	} else {
		line += "├───"
	}

	line += info.Name()

	if !info.IsDir(){
		var size string
		if info.Size() > 0 {
			size = strconv.FormatInt(info.Size(), 10) + "b"
		} else {
			size = "empty"
		}
		line += " (" + size + ")"
	}
	return line
}

func resultTree(out io.Writer, filePath string, printFiles bool, prefix string) error {
	dir, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	list, err := dir.Readdir(-1) //returns a list of directory entries sorted by filename
	dir.Close()
	if err != nil {
		return err
	}

	var res []os.FileInfo
	for idx := range list {
		if list[idx].IsDir() || printFiles {
			res = append(res, list[idx])
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name() < res[j].Name()
	})

	for idx, info := range res {
		isLast := idx == len(res)-1

		if info.Name() != ".DS_Store" {
			fmt.Fprintln(out, printTree(
				isLast,
				prefix,
				info,
			))
		}

		if info.IsDir() {
			var subPrefix string
			if isLast {
				subPrefix = prefix + "\t"
			} else {
				subPrefix = prefix + "│\t"
			}
			err := resultTree(
				out,
				filePath + string(os.PathSeparator) + info.Name(),
				printFiles,
				subPrefix,
			)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}

func dirTree(out io.Writer, filePath string, printFiles bool) error  {
	return resultTree(out, filePath, printFiles, "")
}

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