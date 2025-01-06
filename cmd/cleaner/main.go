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
	daysBeforeStart := flag.Int("days-start", -1, "days before (start)")
	daysBeforeCount := flag.Int("days-count", 1, "days before (count)")
	removeSource := flag.Bool("remove-source", false, "remove source")
	dryRun := flag.Bool("dry-run", false, "dry run")

	flag.Parse()

	if *sourcePath == "" || *destinationPath == "" || *daysBeforeStart <= 0 {
		flag.Usage()
		return
	}

	logFile, err := os.OpenFile("cleaner.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(io.MultiWriter(os.Stderr, logFile), "[cleaner] ", log.LUTC|log.Lmicroseconds|log.Ldate)

	for daysBefore := *daysBeforeStart; daysBefore < *daysBeforeStart+*daysBeforeCount; daysBefore++ {
		err = cleaner.Archive(daysBefore, *sourcePath, *destinationPath, *removeSource, *dryRun)
		if err != nil {
			logger.Fatal(err)
		}
	}
}
