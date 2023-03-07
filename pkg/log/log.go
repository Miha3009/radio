package log

import (
	"log"
	"os"
)

// / -----------------
var (
	infoLogger  *log.Logger = log.Default()
	warnLogger  *log.Logger = log.Default()
	errorLogger *log.Logger = log.Default()
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// ------------------ ^^^^^ remove

type Logger interface {
	Info(string)
	Warn(error)
	Error(error)
	Fatal(err error)
}

func NewLogger() *logger {
	return &logger{}
}

type logger struct{}

func (*logger) Info(message string) {
	infoLogger.Println(message)
}

func (*logger) Warn(err error) {
	warnLogger.Println(err)
}

func (*logger) Error(err error) {
	errorLogger.Println(err)
}

func (*logger) Fatal(err error) {
	errorLogger.Fatalln(err)
}
