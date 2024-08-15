package shared

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(msg string) {
	logger.Println("INFO: " + msg)
}

func Error(msg string) {
	logger.Println("ERROR: " + msg)
}
