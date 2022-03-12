package main

/*
[Duplicate File Handler - Stage 3/4: What's that hash about?](https://hyperskill.org/projects/176/stages/907/implement)
-------------------------------------------------------------------------------
[Public and private scopes](https://hyperskill.org/learn/topic/1894)
[Working with files](https://hyperskill.org/learn/topic/1768)
[Parsing data from strings](https://hyperskill.org/learn/topic/1955)
[Hashing `crypto/md5` package] - PENDING
[Encoding package] - PENDING
*/

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

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

// FileHashMap is a type that uses the hash as a key and a slice of duplicate fileNames as a value
type FileHashMap map[string][]string

func findDuplicateFiles(fileSizes []int, filesMap map[int][]string) (map[int]FileHashMap, error) {
	sameHashMap := make(map[int]FileHashMap)

	for {
		var answer string
		fmt.Println("\nCheck for duplicates?")
		fmt.Scanln(&answer)

		switch answer {
		case "yes":
			for _, fileSize := range fileSizes {
				for _, fileName := range filesMap[fileSize] {
					file, err := os.Open(fileName)
					if err != nil {
						return nil, err
					}

					hash := md5.New()
					_, err = io.Copy(hash, file)
					if err != nil {
						return nil, err
					}
					hashInString := hex.EncodeToString(hash.Sum(nil)[:16])

					if sameHashMap[fileSize] == nil {
						sameHashMap[fileSize] = make(map[string][]string)
					}
					sameHashMap[fileSize][hashInString] = append(sameHashMap[fileSize][hashInString], fileName)

					err = file.Close() // remember to close the file! otherwise, we won't be able to delete
					if err != nil {
						return nil, err
					}
				}
			}
			return sameHashMap, nil

		case "no":
			os.Exit(1) // exit the program if we won't check for duplicates

		default:
			fmt.Println("Wrong option")
		}
	}
}

func getDupFiles(sameHashMap map[int]FileHashMap, fileSizes []int) {
	var counter int

	for _, fileSize := range fileSizes {
		fmt.Println(fileSize, "bytes")
		for hash, files := range sameHashMap[fileSize] {
			if len(files) > 1 {
				fmt.Println("Hash:", hash)
				for i := 0; i < len(files); i++ {
					c := strconv.Itoa(counter + 1)
					// update contents of 'files' to contain the file number like "1. abc.py"
					files[i] = c + ". " + files[i]
					fmt.Printf("%s\n", files[i])
					counter++
				}
				fmt.Println()
			}
		}
	}
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
	rev := getSortingOption()

	// Next we create a map to store the files size, file number and file name
	filesMap := make(map[int][]string)
	// Check for errors and if there is none add the files + their info to the filesMap
	err := addFilesToMap(dir, extension, filesMap)
	if err != nil {
		log.Fatal(err)
	}

	// We create a slice to store the file sizes:
	var fileSizes []int = sortByFileSize(rev, filesMap)

	// Next we need to create the sameHashMap; it will contain all the files that have the same hash.
	sameHashMap := make(map[int]FileHashMap)
	// Check for errors and if there is none add the file size, hash and duplicate files to the sameHashMap
	sameHashMap, err = findDuplicateFiles(fileSizes, filesMap)
	if err != nil {
		log.Fatal(err)
	}

	// Finally, we call getDupFiles; it prints all the files that have the same hash:
	getDupFiles(sameHashMap, fileSizes)
}
