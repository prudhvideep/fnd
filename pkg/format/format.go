package fmt

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/fatih/color"
)

func GetColour(extType string) func(a ...interface{}) string {
	switch extType {
	case "dir":
		return color.New(color.FgHiBlue).SprintFunc()
	case ".txt":
		return color.New(color.FgRed).SprintFunc()
	case ".go":
		return color.New(color.FgGreen).SprintFunc()
	case ".md":
		return color.New(color.FgYellow).SprintFunc()
	case ".mod":
		return color.New(color.FgRed, color.FgHiYellow).SprintFunc()
	case ".exe":
		return color.New(color.FgRed).SprintFunc()
	case ".csv":
		return color.New(color.FgRed).SprintFunc()
	}
	return color.New(color.FgHiCyan).SprintFunc()
}

func FormatOutput(path string, d fs.DirEntry) {

	filename := filepath.Base(path)

	ext := filepath.Ext(filename)

	fileColor := GetColour(ext)
	dirColor := GetColour("dir")

	if d.IsDir() {
		fmt.Println(dirColor(path))
		return
	}

	fmt.Printf("%s \n", filepath.Join(dirColor(filepath.Dir(path)), fileColor(filename)))
}
