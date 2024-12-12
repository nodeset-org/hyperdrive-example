package shared

import (
	"fmt"
	"os"
)

const (
	LogFilePermissions os.FileMode = 0644
)

// Logger for writing to a simple log file
type FileLogger struct {
	// File handle
	file *os.File
}

// Creates a new file logger
func NewFileLogger(path string) (*FileLogger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LogFilePermissions)
	if err != nil {
		return nil, fmt.Errorf("error opening log file [%s]: %w", path, err)
	}

	return &FileLogger{
		file: file,
	}, nil
}

// Writes to the log file
func (f *FileLogger) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

// Closes the log file safely
func (f *FileLogger) Close() error {
	return f.file.Close()
}
