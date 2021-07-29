package logger

import (
	"log"
	"os"
)

type Logger struct {
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
}

var nrdLogger = Logger{
	WarningLogger: log.New(os.Stdout, "NERDCOIN-BOT WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
	InfoLogger: log.New(os.Stdout, "NERDCOIN-BOT INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	ErrorLogger: log.New(os.Stdout, "NERDCOIN-BOT ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
}

func Info(format string)  {
	nrdLogger.InfoLogger.Printf(format)
}

func Warn(format string)  {
	nrdLogger.WarningLogger.Printf(format)
}
func Error(format string) {
	nrdLogger.ErrorLogger.Printf(format)
}
