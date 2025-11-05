package output

import (
	"io"
	"os"
)

// Writer is the global io.Writer you can use in any file
var Writer io.Writer

// Path is the path string used
var Path string

// Init sets up the writer for a given path (file or "stdout")
func Init(path string) error {
	Path = path
	if path == "stdout" {
		Writer = os.Stdout
		return nil
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	Writer = file
	return nil
}
