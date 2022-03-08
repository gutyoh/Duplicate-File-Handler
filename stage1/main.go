package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) == 1 { // check if the only argument is the program file name 'main.go'
		fmt.Println("Directory is not specified")
	} else {
		dir := os.Args[1] // the directory is the second command line argument!
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
			log.Fatal(err)
		}
	}
}
