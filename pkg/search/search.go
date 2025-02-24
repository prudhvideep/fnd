package search

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charlievieth/fastwalk"
	format "github.com/prudhvideep/fnd/pkg/format"
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

func IsHidden(path string) bool {
	parts := strings.Split(filepath.Clean(path), string(os.PathSeparator))
	for _, part := range parts {
		if part[0] == '.' {
			return true
		}
	}
	return false
}

// This method contains the main logic that performs parallel directory traversal
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

	// Get Search Type
	// Valid values - Normal Search, Only Fil
	searchTypeFlag, err := c.PersistentFlags().GetString("type")
	if err != nil {
		log.Fatal("Error reading the type flag")
	}

	searchType := getSearchType(searchTypeFlag)

	//Get extension type
	extType, err := c.PersistentFlags().GetString("ext")
	if err != nil {
		log.Fatal("Error reading the file extension type flag")
	}

	exp, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		log.Fatal("Invalid regex pattern")
	}

	conf := fastwalk.Config{
		Follow:     true,
		NumWorkers: runtime.NumCPU(),
	}

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			return nil
		}

		if IsHidden(path) {
			return nil
		}

		if (searchType == SearchDir && !d.IsDir()) || (searchType == SearchFile && d.IsDir()) {
			return nil
		}

		filePrefix := strings.TrimPrefix(filepath.Ext(d.Name()), ".")

		if extType != "" && extType != filePrefix {
			return nil
		}

		if pattern != "" && exp.MatchString(d.Name()) {
			if relPath, err := filepath.Rel(rootDir, path); err == nil {
				format.FormatOutput(relPath, d)
			}

		}

		return err
	}

	if err := fastwalk.Walk(&conf, rootDir, walkFn); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", rootDir, err)
		os.Exit(1)
	}

}
