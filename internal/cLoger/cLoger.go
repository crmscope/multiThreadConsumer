package cLoger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
	file   *os.File
}

// Конструктор. Создает логер.
func (p *Logger) Init(filePath string, prefix string) error {
	var err error

	p.file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	p.logger = log.New(p.file, prefix, log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

// Метод для ошибок.
func (p *Logger) Error(code int32, message string) {
	p.logger.Println("ERROR [" + fmt.Sprint(code) + "]: " + fmt.Sprintln(message))
}

// Метод для сообщений.
func (p *Logger) Info(message string) {
	p.logger.Println("INFO: " + fmt.Sprintln(message))
}

// Закрывает соединение с файлом.
func (p *Logger) Close() {
	err := p.file.Close()
	if err != nil {
		p.Error(599, "Error closing log file")
	}
}
