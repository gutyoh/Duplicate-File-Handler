package main

/*
[Duplicate File Handler - Stage 3/4: What's that hash about?](https://hyperskill.org/projects/176/stages/907/implement)
-------------------------------------------------------------------------------
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

func sortByFileSize(rev bool, filesMap map[int][]map[int]string) []int {
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
	return fileSizes
}

func findDuplicateFiles(fileSizes []int, filesMap map[int][]map[int]string) map[int]map[string][]string {
	sameHashMap := make(map[int]map[string][]string)

	for {
		var answer string
		fmt.Println("\nCheck for duplicates?")
		fmt.Scanln(&answer)

		if answer == "yes" || answer == "no" {
			if answer == "yes" {
				for _, fileSize := range fileSizes {
					for _, fileNum := range filesMap[fileSize] {
						for _, fileName := range fileNum {
							file, err := os.Open(fileName)
							if err != nil {
								log.Fatal(err)
							}
							// Create a new hash object
							hash := md5.New()
							if _, err = io.Copy(hash, file); err != nil {
								log.Fatal(err)
							}
							hashInString := hex.EncodeToString(hash.Sum(nil)[:16])

							if sameHashMap[fileSize] == nil {
								sameHashMap[fileSize] = make(map[string][]string)
							}
							sameHashMap[fileSize][hashInString] = append(sameHashMap[fileSize][hashInString], fileName)

							err = file.Close() // remember to close the file! otherwise, we won't be able to delete
							if err != nil {
								return nil
							}
						}
					}
				}
			}
			break
		} else {
			os.Exit(1) // exit the program if we won't check for duplicates
		}
	}
	return sameHashMap
}

func getDupFiles(sameHashMap map[int]map[string][]string, fileSizes []int) {
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

	// We create a slice to store the file sizes:
	var fileSizes []int = sortByFileSize(rev, filesMap)

	// Next we need to create the sameHashMap; it will contain all the files that have the same hash.
	var sameHashMap map[int]map[string][]string = findDuplicateFiles(fileSizes, filesMap)

	// Finally, we call getDupFiles; it prints all the files that have the same hash:
	getDupFiles(sameHashMap, fileSizes)
}
