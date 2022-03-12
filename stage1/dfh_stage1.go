package main

/*
[Duplicate File Handler - Stage 1/4: Here come the files](https://hyperskill.org/projects/176/stages/905/implement)
-------------------------------------------------------------------------------
[Primitive types](https://hyperskill.org/learn/topic/1807)
[Input/Output](https://hyperskill.org/learn/topic/1506)
[Slices](https://hyperskill.org/learn/topic/1672)
[Control statements](https://hyperskill.org/learn/topic/1728)
[Errors](https://hyperskill.org/learn/topic/1795)
[Command-line arguments and flags](https://hyperskill.org/learn/topic/1948)
[The `filepath` package] - PENDING
*/

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Directory is not specified")
	} else {
		dir := os.Args[1] // the directory is the second command line argument!
		// use the filepath.Walk function to read 'dir' and print all the files within it:
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
