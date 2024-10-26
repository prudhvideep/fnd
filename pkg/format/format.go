package fmt

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/fatih/color"
)

func FormatOutput(path string, d fs.DirEntry) {
	_, err := d.Info()
	if err != nil {
		fmt.Println(err)
		return
	}

	txtColor := color.New(color.FgRed).SprintFunc()
	goColor := color.New(color.FgGreen).SprintFunc()
	defaultColor := color.New(color.FgHiCyan).SprintFunc()
	dirColor := color.New(color.FgHiBlue).SprintFunc()
	dirColorLight :=  color.New(color.FgBlue).SprintFunc()
	mdColor := color.New(color.FgYellow).SprintFunc()
	modColor := color.New(color.FgRed, color.FgHiYellow).SprintFunc()
	csvColor := color.New(color.FgRed).SprintFunc()

	filename := filepath.Base(path)

	if d.IsDir() {
		fmt.Println(dirColor(path))
		return
	}

	ext := filepath.Ext(filename)
	var coloredFilename string
	switch ext {
	case ".txt":
		coloredFilename = txtColor(filename)
	case ".go":
		coloredFilename = goColor(filename)
	case ".md":
		coloredFilename = mdColor(filename)
	case ".mod":
		coloredFilename = modColor(filename)
	case ".csv":
		coloredFilename = csvColor(filename)
	default:
		coloredFilename = defaultColor(filename)
	}

	fmt.Printf("%s/%s \n", dirColorLight(filepath.Dir(path)), coloredFilename)
}
