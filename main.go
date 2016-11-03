package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
)

func main() {

	if len(os.Args) < 2 {
		color.Red("No search query in arguments")
		os.Exit(0)
	}

	search := os.Args[1]
	fileList := []string{}

	err := filepath.Walk("/var/log", func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	hits := 0

	var wg sync.WaitGroup
	wg.Add(len(fileList))

	color.Green("Term: " + search)
	color.Green("Searching trough : " + strconv.Itoa(len(fileList)) + " files")

	for _, path := range fileList {

		go func(path string, search string) {

			file, err := ioutil.ReadFile(path)

			if err != nil {
				color.Red(err.Error())
			}

			lines := strings.Split(string(file), "\n")

			for line, content := range lines {
				bingo := strings.Contains(strings.ToLower(content), strings.ToLower(search))

				if bingo {
					color.Yellow(path + " (" + "LINE " + strconv.Itoa(line) + ") ")
					color.Cyan(content)
					hits = hits + 1
				}
			}

			wg.Done()

		}(path, search)

	}

	wg.Wait()
}
