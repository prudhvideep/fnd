package search

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/prudhvideep/fnd/pkg/config"
	format "github.com/prudhvideep/fnd/pkg/format"
	mySsh "github.com/prudhvideep/fnd/pkg/ssh"
	ssh "golang.org/x/crypto/ssh"
)

type PrintMessage struct {
	path string
	d    fs.DirEntry
}

func RemoteSearch(args []string, typeFlag string, directory string, credntials *config.Credentials) error {
	fmt.Println("Inside the Remote Search")

	var pattern string

	if len(args) == 0 {
		pattern = "*"
	}

	if len(args) > 0 {
		pattern = args[0]

		if _, err := regexp.Compile(pattern); err != nil {
			pattern = "*" + pattern + "*"
		}
	}

	client, err := mySsh.InitializeConnection(credntials)
	if err != nil {
		return err
	}

	output, err := RunFindCommand(client, pattern, directory)
	if err != nil {
		return err
	}

	fmt.Println(output)

	return nil
}

func RunFindCommand(client *ssh.Client, pattern, directory string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var cmd string
	if pattern == "" {
		pattern = "*"
	}
	if directory == "" {
		directory = "."
	}

	cmd = fmt.Sprintf("find %s -name '%s'", directory, pattern)

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(cmd)
	if err != nil {
		return stderr.String(), err
	}

	// err = session.Run("echo $home")
	// if err != nil{
	// 	return stderr.String(), err
	// }

	// fmt.Println(stdout.String())

	return stdout.String(), nil
}

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
	const workers = 10
	var wg sync.WaitGroup
	pathChannel := make(chan PrintMessage)

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

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for pmsg := range pathChannel {
				format.FormatOutput(pmsg.path,pmsg.d)
			}
		}()
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
			pathChannel <- PrintMessage{
				path: path,
				d:    d,
			}
		}
		return nil
	})

	close(pathChannel)

	wg.Wait()

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}

func fetchDir(args []string, directory string) {
	var pattern string
	var isRegExp bool
	var wg sync.WaitGroup
	const workers = 10
	pathChannel := make(chan string)

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

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for path := range pathChannel {
				fmt.Println(path)
			}
		}()
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
				pathChannel <- path
			}
		}
		return nil
	})

	wg.Done()

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}

func fetchFiles(args []string, directory string) {
	var pattern string
	var isRegExp bool
	var wg sync.WaitGroup
	const workers = 10
	pathChannel := make(chan PrintMessage)

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

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for pmsg := range pathChannel {
				format.FormatOutput(pmsg.path, pmsg.d)
			}
		}()
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
				pathChannel <- PrintMessage{
					path: path,
					d:    d,
				}
			}
		}
		return nil
	})

	close(pathChannel)
	wg.Wait()

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}
