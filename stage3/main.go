package main

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

		dir := os.Args[1] // the directory is the second command line argument!
		dir = strings.Join(os.Args[1:], " ")

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

		// Create a slice to sort the filesMap:
		keys := make([]int, 0, len(filesMap))
		for k := range filesMap {
			keys = append(keys, k)
		}

		// If 'rev' is true, then sort the slice in descending order based on SIZE of the files
		// Otherwise, sort in ascending order:
		if rev {
			sort.Sort(sort.Reverse(sort.IntSlice(keys)))
		} else {
			sort.Sort(sort.IntSlice(keys))
		}

		// Print the size in bytes and afterwards the sorted files:
		for _, k := range keys {
			fmt.Println(k, "bytes")
			for _, v := range filesMap[k] {
				fmt.Println(v)
			}
			fmt.Println()
		}

		for {
			var answer string
			fmt.Println("Check for duplicates?")
			fmt.Scanln(&answer)

			var sameHashMap = map[int]map[string][]string{} // create a map to store the files with the same size

			if answer == "yes" || answer == "no" {
				if answer == "yes" {
					// create a md5 hash with md5.New() and calculate the hash for each file
					// and store the file name and the hash in a map:
					for _, k := range keys {
						for _, v := range filesMap[k] {
							file, err := os.Open(v)
							if err != nil {
								log.Fatal(err)
							}
							defer file.Close()

							hash := md5.New()
							if _, err := io.Copy(hash, file); err != nil {
								log.Fatal(err)
							}
							hashInBytes := hash.Sum(nil)[:16]
							hashInString := hex.EncodeToString(hashInBytes)

							if sameHashMap[k] == nil {
								sameHashMap[k] = make(map[string][]string)
							}
							sameHashMap[k][hashInString] = append(sameHashMap[k][hashInString], v)
						}
					}
					// iterate over the map and only print the files that have the same hash:
					counter := 0
					i := 0
					for _, k := range keys {
						fmt.Println(k, "bytes")
						for h, v := range sameHashMap[k] {
							if len(v) > 1 {
								fmt.Println("Hash:", h)
								for _, v2 := range v {
									c := strconv.Itoa(counter + 1)
									// update contents of 'v' to become c + ". " + v
									v[i] = c + ". " + v[i]
									v2 = c + ". " + v2
									fmt.Printf("%s\n", v2)
									counter++
									i++
								}
								i = 0
								fmt.Println()
							}
						}
					}
					break
				} else {
					break
				}
			}
		}
	}
}
