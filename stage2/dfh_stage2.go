package main

/*
[Duplicate File Handler - Stage 2/4: How much does it weigh?](https://hyperskill.org/projects/176/stages/906/implement)
-------------------------------------------------------------------------------
[Loops](https://hyperskill.org/learn/topic/1531)
[Public and private scopes](https://hyperskill.org/learn/topic/1894)
[Functions](https://hyperskill.org/learn/topic/1750)
[Function decomposition](https://hyperskill.org/learn/topic/1893)
[Maps](https://hyperskill.org/learn/topic/1824)
[Operations with maps](https://hyperskill.org/learn/topic/1850)
[Sorting slices] - PENDING UPLOAD TO HYPERSKILL
*/

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var rev bool // global variable 'rev' to determine the sorting order based on SIZE of the files

func getArgs(args []string) []string {
	args = os.Args

	if len(args) == 1 {
		fmt.Println("Directory is not specified")
		os.Exit(1)
	}
	return args
}

func getExtension() string {
	var extension string

	fmt.Println("Enter file format:")
	fmt.Scanln(&extension)

	if len(extension) == 0 {
		return ""
	} else {
		return "." + extension
	}
}

func getSortingOption() bool {
	var n int
	fmt.Println("Size sorting options:\n1. Ascending\n2. Descending")

	for {
		fmt.Println("Enter a sorting option:")
		fmt.Scanln(&n)

		if n == 1 || n == 2 {
			if n == 1 {
				rev = true
			} else {
				rev = false
			}
			break
		} else {
			fmt.Println("Wrong option")
		}
	}
	return rev
}

func addFilesToMap(dir string, extension string, filesMap map[int][]map[int]string) {
	fileNum := 1

	// if the extension is NOT specified, then add all the files to the map
	if len(extension) == 0 {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if info.IsDir() {
				return nil
			}
			filesMap[int(info.Size())] = append(filesMap[int(info.Size())], map[int]string{
				fileNum: path,
			})
			fileNum++
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

	} else { // if the extension is specified, then add only the files with the specified extension to the map
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) == extension {
				filesMap[int(info.Size())] = append(filesMap[int(info.Size())], map[int]string{
					fileNum: path,
				})
				fileNum++
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func sortByFileSize(rev bool, filesMap map[int][]map[int]string) {
	fileSizes := make([]int, 0, len(filesMap))

	for fileSize := range filesMap {
		fileSizes = append(fileSizes, fileSize)
	}

	if rev {
		sort.Sort(sort.Reverse(sort.IntSlice(fileSizes)))
	} else {
		sort.Sort(sort.IntSlice(fileSizes))
	}

	// Print the sorted sizes in bytes and afterwards the respective file names:
	for _, fileSize := range fileSizes {
		fmt.Println()
		fmt.Println(fileSize, "bytes")
		for _, fileNum := range filesMap[fileSize] {
			for _, fileName := range fileNum {
				fmt.Println(fileName)
			}
		}
	}
}

func main() {
	// The first step is to get the command-line arguments passed to our program:
	getArgs(os.Args)
	// Since the directory is the second command line argument, we create 'dir' and store it there:
	dir := strings.Join(os.Args[1:], " ")

	// Take as an input the extension of files to check; if nothing is entered we check all files in the dir.
	extension := getExtension()

	// Take as an input the sorting option - 1 for ascending; 2 for descending.
	rev = getSortingOption()

	// Next we create a map to store the files size, file number and file name
	filesMap := make(map[int][]map[int]string)
	addFilesToMap(dir, extension, filesMap) // we add the files + their info to the map

	// We call the sortByFileSize function to sort files by size and print them out:
	sortByFileSize(rev, filesMap)
}
