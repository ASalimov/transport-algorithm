package main

import (
	"log"
	"os"
)

// Info is logger to stdout
var Info *log.Logger

// Error is logger to stderr
var Error *log.Logger

func init() {
	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
