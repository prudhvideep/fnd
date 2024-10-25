package search

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
)

func Find(args []string, typeFlag string, directory string) {
	switch typeFlag {
	case "d":
		fetchDir(args, directory)
	case "f":
		fetchFiles(args, directory)
	default:
		fetchAll(args, directory)
	}
}

func fetchAll(args []string, directory string) {

	var pattern string
	var isRegExp bool

	if len(args) == 0 {
		pattern = "*"
	}

	if len(args) > 0 {
		pattern = args[0]

		if _, err := regexp.Compile(pattern); err == nil {
			isRegExp = true
		} else {
			pattern = "*" + pattern + "*"
		}
	}

	if directory == "" {
		directory = "."
	}

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if filepath.Base(path)[0] == '.' {
			return nil
		}

		var match bool

		if isRegExp {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return err
			}
			match = re.MatchString(d.Name())
		} else {
			var e error
			match, e = filepath.Match(pattern, d.Name())
			if e != nil {
				return e
			}
		}

		if match {
			fmt.Println(path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}

func fetchDir(args []string, directory string) {
	var pattern string
	var isRegExp bool

	if len(args) == 0 {
		pattern = "*"
	}

	if len(args) > 0 {
		pattern = args[0]

		if _, err := regexp.Compile(pattern); err == nil {
			isRegExp = true
		} else {
			pattern = "*" + pattern + "*"
		}
	}

	if directory == "" {
		directory = "."
	}

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if filepath.Base(path)[0] == '.' {
			return nil
		}

		var match bool

		if d.IsDir() {
			if isRegExp {
				re, err := regexp.Compile(pattern)
				if err != nil {
					return err
				}
				match = re.MatchString(d.Name())
			} else {
				var e error
				match, e = filepath.Match(pattern, d.Name())
				if e != nil {
					return e
				}
			}

			if match {
				fmt.Println(path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}

func fetchFiles(args []string, directory string) {
	var pattern string
	var isRegExp bool

	if len(args) == 0 {
		pattern = "*"
	}

	if len(args) > 0 {
		pattern = args[0]

		if _, err := regexp.Compile(pattern); err == nil {
			isRegExp = true
		} else {
			pattern = "*" + pattern + "*"
		}
	}

	if directory == "" {
		directory = "."
	}

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if filepath.Base(path)[0] == '.' {
			return nil
		}

		var match bool

		if !d.IsDir() {
			if isRegExp {
				re, err := regexp.Compile(pattern)
				if err != nil {
					return err
				}
				match = re.MatchString(d.Name())
			} else {
				var e error
				match, e = filepath.Match(pattern, d.Name())
				if e != nil {
					return e
				}
			}

			if match {
				fmt.Println(path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}
