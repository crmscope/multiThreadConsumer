package cLoger

import (
	"log"
	"os"
)

type FileLogger struct {
	logger *log.Logger
	file   *os.File
}

// NewFileLogger creates a new instance of the logger that writes to the specified file.
func NewFileLogger(filePath string, prefix string) (*FileLogger, error) {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	logger := log.New(logFile, prefix, log.Ldate|log.Ltime|log.Lshortfile)
	return &FileLogger{
		logger: logger,
		file:   logFile,
	}, nil
}

// Info writes an informational message
func (fl *FileLogger) Info(msg string) {
	fl.logger.Println("INFO: " + msg)
}

// Error writes an error message
func (fl *FileLogger) Error(msg string) {
	fl.logger.Println("ERROR: " + msg)
}

// Close closes the log file
func (fl *FileLogger) Close() error {
	return fl.file.Close()
}
