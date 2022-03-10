package main

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

// global "SET" possibleAns to check for "yes" or "no" input ONLY.
var possibleAns = map[string]bool{
	"yes": true,
	"no":  false,
}

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
		// fmt.Println("The EXTENSION is:", extension)
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

func checkDuplicateFiles(fileSizes []int, filesMap map[int][]map[int]string) map[int]map[string][]string {
	sameHashMap := make(map[int]map[string][]string)

	for {
		var answer string
		fmt.Println("\nCheck for duplicates?")
		fmt.Scanln(&answer)

		// check if answer is in possibleAns map
		if _, ok := possibleAns[answer]; ok {
			if possibleAns[answer] {
				for _, fileSize := range fileSizes {
					for _, fileNum := range filesMap[fileSize] {
						for _, fileName := range fileNum {
							file, err := os.Open(fileName)
							if err != nil {
								log.Fatal(err)
							}

							hash := md5.New()
							if _, err := io.Copy(hash, file); err != nil {
								log.Fatal(err)
							}
							hashInBytes := hash.Sum(nil)[:16]
							hashInString := hex.EncodeToString(hashInBytes)

							if sameHashMap[fileSize] == nil {
								sameHashMap[fileSize] = make(map[string][]string)
							}
							sameHashMap[fileSize][hashInString] = append(sameHashMap[fileSize][hashInString], fileName)

							err = file.Close()
							if err != nil {
								return nil
							}
						}
					}
				}
			}
			break
		} else {
			// exit the program if we won't check for duplicates
			os.Exit(1)
		}
	}
	return sameHashMap
}

func getFilesSameHash(sameHashMap map[int]map[string][]string, fileSizes []int) []int {
	var fileNums []int
	counter := 0

	for _, fileSize := range fileSizes {
		fmt.Println(fileSize, "bytes")
		for hash, files := range sameHashMap[fileSize] {
			if len(files) > 1 {
				fmt.Println("Hash:", hash)
				for i := 0; i < len(files); i++ {
					c := strconv.Itoa(counter + 1)
					// update contents of 'files' to become c + ". " + files
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

func deleteDuplicates(sameHashMap map[int]map[string][]string, fileSizes []int, fileNums []int) {
	deletedFileSize := 0
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

	// sort filesToDelete
	sort.Ints(filesToDelete)

	if contains(filesToDelete, fileNums) {
		cnt := 0
		for _, fileSize := range fileSizes {
			for _, files := range sameHashMap[fileSize] {
				if len(files) > 1 {
					for i := 0; i < len(files); i++ {
						// add to deletedFileSize the size of the file that is being deleted:
						if cnt == len(filesToDelete) {
							break
						}

						// get the file number by using the split function
						fileNum := strings.Split(files[i], ".")

						// if fileNum is the same as the filesToDelete then delete the file:
						if fileNum[0] == strconv.Itoa(filesToDelete[cnt]) {
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
							continue
						}

						cnt++
					}
				}
			}
		}
	}
	fmt.Println("Total freed up space:", deletedFileSize, "bytes")
	os.Exit(1)
}

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
	// getArgs checks if the only argument is the program name
	getArgs(os.Args)

	// next we make the user enter the extension of the files we want to check
	extension := getExtension()

	// next we make the user enter the sorting option
	rev = getSortingOption()

	// since the directory is the second command line argument, we create 'dir' and store it there:
	dir := strings.Join(os.Args[1:], " ")

	// next we create a map to store the files size, file number and file name
	filesMap := make(map[int][]map[int]string)
	addFilesToMap(dir, extension, filesMap)

	// create a slice to store the file sizes
	fileSizes := sortByFileSize(rev, filesMap)

	// next we need to check for duplicates and create the md5 hash for each file
	sameHashMap := checkDuplicateFiles(fileSizes, filesMap)

	// now we need to get the files with the same hash
	fileNums := getFilesSameHash(sameHashMap, fileSizes)

	// Finally, we delete the duplicate files
	deleteDuplicates(sameHashMap, fileSizes, fileNums)

}
