package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Directory is not specified")
	} else {
		dir := strings.Join(os.Args[1:], " ") // the directory is the second command line argument!
		// use the 'Walk' function to read 'dir' and print all the files
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			fmt.Println(path)
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}
