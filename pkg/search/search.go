package search

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	regexp "github.com/wasilibs/go-re2"
)

const (
	SearchNormal = iota
	SearchDir
	SearchFile
)

const (
	CAPACITY = 1000
)

type Entry struct {
	Path  string
	Name  string
	IsDir bool
}

func (e *Entry) IsHidden() bool {
	entryName := e.Name

	if strings.HasPrefix(entryName, ".") {
		return true
	}

	if runtime.GOOS == "windows" {
		pointer, err := syscall.UTF16PtrFromString(entryName)
		if err != nil {
			return false
		}
		attributes, err := syscall.GetFileAttributes(pointer)
		if err != nil {
			return false
		}

		return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
	}

	return false
}

func getSearchType(searchType string) int {
	switch searchType {
	case "f":
		return SearchFile
	case "d":
		return SearchDir
	default:
		return SearchNormal
	}
}

//This method contains the main logic that performs parallel directory traversal
func Find(c *cobra.Command, args []string) {
	if len(args) > 2 {
		log.Fatal("Only two arguments are allowed")
	}

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error reading the parent dir")
	}

	if len(args) == 2 {
		reqSearchDir := args[1]
		rootDir = filepath.Join(rootDir, reqSearchDir)
	}

	pattern := ".*"
	if len(args) >= 1 {
		pattern = args[0]
	}

	//Get Search Type
	//Valid values - Normal Search, Only Fil
	typeFlag, err := c.PersistentFlags().GetString("type")
	if err != nil {
		log.Fatal("Error reading the type flag")
	}

	searchType := getSearchType(typeFlag)

	exp, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal("Invalid regex pattern")
	}

	entryChan := make(chan Entry, CAPACITY)
	resultChan := make(chan Entry, CAPACITY)
	var wg sync.WaitGroup

	workers := runtime.NumCPU() * 2
	wg.Add(workers)
	for i := 1; i <= workers; i++ {
		go ProcessEntries(entryChan, resultChan, pattern, exp, searchType, &wg)
	}

	go func() {
		Walk(rootDir, entryChan)
		close(entryChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for {
		res, ok := <-resultChan
		if !ok {
			return
		}

		cleanRoot := filepath.Clean(rootDir)
		cleanRes := filepath.Clean(res.Path)
		relPath, err := filepath.Rel(cleanRoot, cleanRes)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println(relPath)
	}
}

func Walk(root string, entryChan chan<- Entry) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}

	for _, entry := range entries {

		if entry.Name()[0] == '.' {
			continue
		}

		entryChan <- Entry{Name: entry.Name(), Path: filepath.Join(root, entry.Name()), IsDir: entry.IsDir()}

		if entry.IsDir() {
			Walk(filepath.Join(root, entry.Name()), entryChan)
		}
	}
}

func ProcessEntries(entryChan <-chan Entry, resultChan chan<- Entry, pattern string, exp *regexp.Regexp, searchType int, wg *sync.WaitGroup) {

	defer wg.Done()

	for entry := range entryChan {

		if searchType == SearchNormal || (searchType == SearchDir && entry.IsDir) || (searchType == SearchFile && !entry.IsDir) {
			if exp.MatchString(entry.Name) {
				resultChan <- entry
			}
		}
	}

}
