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

func getArgs(args []string) []string {
	args = os.Args
	if len(args) == 1 {
		fmt.Println("Directory is not specified")
		os.Exit(1)
	}
	return args
}

func getExtension(extension string) string {
	fmt.Println("Enter file format:")
	fmt.Scanln(&extension)

	if len(extension) == 0 {
		return ""
	} else {
		return "." + extension
	}
}

// getSortingOption returns the sorting option based on two options '1': desc order and '2': asc order
func getSortingOption(n int) bool {
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
		}
		return rev
	}
}

func addFilesToMap(dir string, extension string, filesMap map[int][]string) {
	// If the extension is "" (not specified) then add all files to the map
	if len(extension) == 0 {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if info.IsDir() {
				return nil
			}
			filesMap[int(info.Size())] = append(filesMap[int(info.Size())], path)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	} else { // If the extension is specified, then add only the files with the specified extension to the map
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) == "."+extension {
				filesMap[int(info.Size())] = append(filesMap[int(info.Size())], path)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func sortFilesMap(filesMap map[int][]string) []int {
	var keys []int
	for k := range filesMap {
		keys = append(keys, k)
	}
	if rev {
		sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	} else {
		sort.Ints(keys)
	}
	return keys
}

// func checkDuplicateFiles(filesMap map[int][]string, keys []int) {
// 	var answer string
// 	fmt.Println("Check for duplicates?")
// 	fmt.Scanln(&answer)

// 	var sameHashMap = map[int]map[string][]string{} // create a new map to store the files with the same size

// 	if answer == "yes" || answer == "no" {
// 		if answer == "yes" {
// 			// create an md5 hash with the md5.New() func and calculate the hash for each file
// 			// then store the file name and the hash in a map:
// 		}
// 	}
// }

func main() {
	if len(os.Args) == 1 { // check if the only argument is the program file name 'main.go'
		fmt.Println("Directory is not specified")
	} else {
		var extension string
		fmt.Println("Enter file format:")
		fmt.Scanln(&extension)

		fmt.Println("Size sorting options:\n1. Ascending\n2. Descending")

		for {
			var n int
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

		filesMap := make(map[int][]string) // create a map to store the file names and their sizes

		dir := strings.Join(os.Args[1:], " ") // the directory is the second command line argument!

		// if the extension is NOT specified, then add all the files to the map
		if len(extension) == 0 {
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Fatal(err)
				}
				if info.IsDir() {
					return nil
				}
				filesMap[int(info.Size())] = append(filesMap[int(info.Size())], path)
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
				if filepath.Ext(path) == "."+extension {
					// fmt.Println("Reading file:", path, "with extension:", extension)
					filesMap[int(info.Size())] = append(filesMap[int(info.Size())], path)
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		}

		// Create a slice to store the file sizes
		fileSizes := make([]int, 0, len(filesMap))
		for fileSize := range filesMap {
			fileSizes = append(fileSizes, fileSize)
		}

		// If 'rev' is true, then sort the fileSizes slice in descending order:
		// Otherwise, sort in ascending order:
		if rev {
			sort.Sort(sort.Reverse(sort.IntSlice(fileSizes)))
		} else {
			sort.Sort(sort.IntSlice(fileSizes))
		}

		// Print the sorted sizes in bytes and afterwards the respective file names:
		for _, fileSize := range fileSizes {
			fmt.Println(fileSize, "bytes")
			for _, fileName := range filesMap[fileSize] {
				fmt.Println(fileName)
			}
			fmt.Println()
		}

		// For loop to check if there are any duplicate files:
		for {
			var answer string
			fmt.Println("Check for duplicates?")
			fmt.Scanln(&answer)

			var sameHashMap = map[int]map[string][]string{} // create a map to store the files with the same size
			// var sameHashMap = map[int]map[string]map[int][]string{} // create a map to store the files with the same size

			if answer == "yes" || answer == "no" {
				if answer == "yes" {
					// create a md5 hash with md5.New() and calculate the hash for each file
					// and store the file name and the hash in a map:
					for _, fileSize := range fileSizes {
						for _, fileName := range filesMap[fileSize] {
							file, err := os.Open(fileName)
							if err != nil {
								log.Fatal(err)
							}

							// Create a new hash object with the md5.New() func:
							hash := md5.New()
							if _, err := io.Copy(hash, file); err != nil {
								log.Fatal(err)
							}
							hashInBytes := hash.Sum(nil)[:16]
							hashInString := hex.EncodeToString(hashInBytes) // convert the hash to a string

							if sameHashMap[fileSize] == nil {
								sameHashMap[fileSize] = make(map[string][]string)
							}
							sameHashMap[fileSize][hashInString] = append(sameHashMap[fileSize][hashInString], fileName)

							err = file.Close()
							if err != nil {
								return
							}
						}
					}
					// iterate over the map and only print the files that have the same hash:
					counter := 0
					// i := 0
					var fileNums []int
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
									// i++
								}
								// i = 0
								fmt.Println()
							}
						}
					}

					deletedFileSize := 0 // create a variable to accumulate the sum of the size of deleted files
					var filesToDelete []int

					fmt.Println("Delete files?")
					fmt.Scanln(&answer)
					if answer == "yes" || answer == "no" {
						if answer == "yes" {
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
					}
				}
			}
		}
	}
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
