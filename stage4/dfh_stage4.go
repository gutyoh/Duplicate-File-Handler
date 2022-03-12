package main

/*
[Duplicate File Handler - Stage 4/4: Delete them all](https://hyperskill.org/projects/176/stages/908/implement)
-------------------------------------------------------------------------------
[Advanced input](https://hyperskill.org/learn/topic/2027)
[Operations with strings](https://hyperskill.org/learn/topic/2023)
*/

import (
	"bufio"
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

func findDuplicateFiles(fileSizes []int, filesMap map[int][]string) map[int]map[string][]string {
	sameHashMap := make(map[int]map[string][]string)

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
						log.Fatal(err)
					}

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
			return sameHashMap

		case "no":
			os.Exit(1) // exit the program if we won't check for duplicates

		default:
			fmt.Println("Wrong option")
		}
	}
}

// We update the previous getDupFiles function to: getDupFileNums
// it returns an array that contains all the duplicate file numbers
func getDupFileNums(sameHashMap map[int]map[string][]string, fileSizes []int) []int {
	var fileNums []int
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

					fileNums = append(fileNums, counter+1) // store the duplicate file numbers
					counter++
				}
				fmt.Println()
			}
		}
	}
	return fileNums
}

func readDupFileNums() []int {
	var filesToDelete []int

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter file numbers to delete:")
	scanner.Scan()
	line := scanner.Text()

	for len(line) != 0 || line == "" {
		splitLine := strings.Split(line, " ")

		if len(filesToDelete) >= len(splitLine) {
			break
		}
		for _, fileNum := range splitLine {
			// if num is an integer then append it to the slice
			if num, err := strconv.Atoi(fileNum); err == nil {
				filesToDelete = append(filesToDelete, num)
			} else {
				fmt.Println("Wrong format")
				fmt.Println("Enter file numbers to delete:")
				scanner.Scan()
				line = scanner.Text()
			}
		}
	}
	sort.Ints(filesToDelete) // remember to sort files in ascending order before returning!
	return filesToDelete
}

func deleteDupFiles(sameHashMap map[int]map[string][]string, fileSizes []int, dupFileNums []int, filesToDelete []int) {
	var deletedFileSize int

	if contains(filesToDelete, dupFileNums) {
		var counter int
		for _, fileSize := range fileSizes {
			for _, files := range sameHashMap[fileSize] {
				if len(files) > 1 {
					for i := 0; i < len(files); i++ {
						// add to deletedFileSize the size of the file that is being deleted:
						if counter == len(filesToDelete) {
							break
						}
						// get the file number by using the split function
						fileNum := strings.Split(files[i], ".")

						// if fileNum is the same as the filesToDelete then delete the file:
						if fileNum[0] == strconv.Itoa(filesToDelete[counter]) {
							// to delete the file remove the prefix 1.:
							fileName := strings.TrimPrefix(files[i], fileNum[0]+". ")

							// get the file size in bytes of 'fileName':
							fileInfo, err := os.Stat(fileName)
							if err != nil {
								fmt.Println(err)
							}
							deletedFileSize += int(fileInfo.Size())

							err = os.Remove(fileName)
							if err != nil {
								fmt.Println(err)
							}
						} else {
							break
						}
						counter++
					}
				}
			}
		}
	}
	fmt.Println("Total freed up space:", deletedFileSize, "bytes")
	os.Exit(1)
}

// Contains is a function to help us validate if the files we read from the input to be deleted
// are actually within the fileNums slice
func contains(s []int, e []int) bool {
	for _, a := range s {
		for _, b := range e {
			if a == b {
				return true
			}
		}
	}
	return false
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

	// We create a slice to store the file sizes:
	var fileSizes []int = sortByFileSize(rev, filesMap)

	// Next we need to create the sameHashMap; it will contain all the files that have the same hash.
	var sameHashMap map[int]map[string][]string = findDuplicateFiles(fileSizes, filesMap)

	// We create a slice to contain the number of the files that have the same hash (duplicates)
	var dupFileNums []int = getDupFileNums(sameHashMap, fileSizes)

	// We take as an input the file numbers to delete and store them in the filesToDelete slice
	var filesToDelete []int = readDupFileNums()

	// Finally, we delete the files that have the same hash (duplicates)
	deleteDupFiles(sameHashMap, fileSizes, dupFileNums, filesToDelete)
}
