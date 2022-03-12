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
[Sorting slices](https://hyperskill.org/learn/topic/2010)
*/

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var rev bool // global variable 'rev' to determine the sorting order based on SIZE of the files

func getExtension() string {
	var extension string

	fmt.Println("Enter file format:")
	fmt.Scanln(&extension)

	if extension == "" {
		return ""
	}
	return "." + extension
}

func getSortingOption() bool {
	var n int
	fmt.Println("Size sorting options:\n1. Ascending\n2. Descending")

	for {
		fmt.Scanln(&n)
		switch n {
		case 1:
			return true
		case 2:
			return false
		default:
			fmt.Println("Wrong option")
		}
	}
}

func addFilesToMap(dir string, extension string, filesMap map[int][]string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if extension == "" || filepath.Ext(path) == extension {
			filesMap[int(info.Size())] = append(filesMap[int(info.Size())], path)
		}
		return nil
	})
}

func sortByFileSize(rev bool, filesMap map[int][]string) []int {
	fileSizes := make([]int, 0, len(filesMap))

	for fileSize := range filesMap {
		fileSizes = append(fileSizes, fileSize)
	}

	if rev {
		sort.Sort(sort.Reverse(sort.IntSlice(fileSizes)))
	} else {
		sort.Ints(fileSizes)
	}

	for _, fileSize := range fileSizes {
		fmt.Println()
		fmt.Println(fileSize, "bytes")
		for _, fileName := range filesMap[fileSize] {
			fmt.Println(fileName)
		}
	}

	return fileSizes
}

func main() {
	// The first step is to get the command-line arguments passed to our program:
	if len(os.Args) == 1 {
		fmt.Println("Directory is not specified")
		return
	}
	// Since the directory is the second command line argument, we create 'dir' and store it there:
	dir := os.Args[1]

	// Take as an input the extension of files to check; if nothing is entered we check all files in the dir.
	extension := getExtension()

	// Take as an input the sorting option - 1 for ascending; 2 for descending.
	rev = getSortingOption()

	// Next we create a map to store the files size, file number and file name
	filesMap := make(map[int][]string)
	// Check for errors and if there is none add the files + their info to the map
	err := addFilesToMap(dir, extension, filesMap)
	if err != nil {
		log.Fatal(err)
	}

	// We call the sortByFileSize function to sort files by size and print them out:
	sortByFileSize(rev, filesMap)
}
