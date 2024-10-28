package fmt

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func RgbColor(text string, r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, text)
}

func FormatOutput(path string, d fs.DirEntry) {
	_, err := d.Info()
	if err != nil {
		fmt.Println(err)
		return
	}

	filename := filepath.Base(path)

	if d.IsDir() {
		fmt.Printf("%s \n", RgbColor(path, 86, 299, 245))
		return
	}

	ext := filepath.Ext(filename)
	var coloredFilename string
	switch ext {
	case ".txt":
		coloredFilename = RgbColor(filename, 245, 85, 64)
	case ".go":
		coloredFilename = RgbColor(filename, 48, 242, 58)
	case ".md":
		coloredFilename = RgbColor(filename, 255, 165, 0)
	case ".mod":
		coloredFilename = RgbColor(filename, 235, 137, 101)
	case ".csv":
		coloredFilename = RgbColor(filename, 245, 85, 64)
	case ".py":
		coloredFilename = RgbColor(filename, 48, 242, 58)
	case ".java":
		coloredFilename = RgbColor(filename, 48, 242, 58)
	default:
		coloredFilename = RgbColor(filename, 241, 245, 122)
	}

	if filepath.Dir(path) == "." {
		fmt.Printf("%s \n", coloredFilename)
	}else{
		fmt.Printf("%s/%s \n", RgbColor(filepath.Dir(path), 86, 299, 245), coloredFilename)
	}
}
