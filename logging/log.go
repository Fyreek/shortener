package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logLevel = Debug
var logToFile = false

// Init initializes logging
func Init(absPath, fileLogging bool, path string) error {
	ex, err := os.Executable()
	if err != nil {
		Log(Failure, "Executable path could not be read for logging")
		return err
	}
	if !absPath {
		path = filepath.Dir(ex) + "/" + path
	}
	logToFile = fileLogging
	if logToFile {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		// TODO: Add proper file closing later
		// defer file.Close()
		log.SetOutput(file)
		return nil
	}
	return nil
}

// Log logs output to the standard console window
func Log(lType LogType, a ...interface{}) {
	if lType < logLevel {
		return
	}
	lTime := "[" + time.Now().Format(time.RFC3339) + "]"
	lLevel := ""
	switch lType {
	case Debug:
		lLevel = "[Debug]"
	case Info:
		lLevel = "[Info]"
	case Failure:
		lLevel = "[Failure]"
	default:
		lLevel = "[Debug]"
	}
	fmt.Println(lTime, lLevel, a)
	if logToFile {
		log.Println(lLevel, a)
	}
}