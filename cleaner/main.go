package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/pkierski/wokanda-scrapper/pkg/cleaner"
)

func main() {
	sourcePath := flag.String("source", "storage", "source path")
	destinationPath := flag.String("destination", "storage", "destination path")
	daysBefore := flag.Int("days", 7, "days before")
	removeSource := flag.Bool("remove-source", false, "remove source")
	flag.Parse()

	logFile, err := os.OpenFile("cleaner.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(io.MultiWriter(os.Stderr, logFile), "[cleaner] ", log.LUTC|log.Lmicroseconds|log.Ldate)

	err = cleaner.Archive(*daysBefore, *sourcePath, *destinationPath, *removeSource)
	if err != nil {
		logger.Fatal(err)
	}
}
