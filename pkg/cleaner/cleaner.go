package cleaner

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func createArchive(files []string, buf io.Writer) error {
	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	// gw := gzip.NewWriter(buf)
	gw, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Iterate over files and add them to the tar archive
	for _, file := range files {
		err := addToArchive(tw, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func addToArchive(tw *tar.Writer, filename string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = filename

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}

func archPathAndFiles(date time.Time, daysBefore int, filesPath, archPath string) (string, []string) {
	date = date.AddDate(0, 0, -daysBefore)
	pattern := fmt.Sprintf("trials_%v_*.json", date.Format("2006-01-02"))
	archName := fmt.Sprintf("trials_%v.tar.gz", date.Format("2006-01-02"))

	files, err := filepath.Glob(filepath.Join(filesPath, pattern))
	if err != nil {
		log.Fatalf("Error getting files from glob: %v", err)
	}
	return filepath.Join(archPath, archName), files
}

func archiveFiles(output string, files []string) error {
	if _, err := os.Stat(output); err == nil {
		return fmt.Errorf("archive already exists: %v", output)
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	return createArchive(files, out)
}

func Archive(daysBefore int, filesPath, archPath string, removeSource bool, dryRun bool) error {
	archName, files := archPathAndFiles(time.Now(), daysBefore, filesPath, archPath)

	if dryRun {
		fmt.Printf("Would archive %v files to %v\n", len(files), archName)
		return nil
	}

	if len(files) == 0 {
		return nil
	}

	err := os.MkdirAll(archPath, 0o644)
	if err != nil {
		return err
	}

	err = archiveFiles(archName, files)
	if err != nil {
		return err
	}

	var errs []error
	if removeSource {
		errs = make([]error, len(files))
		for i, file := range files {
			errs[i] = os.Remove(file)
		}
	}

	return errors.Join(errs...)
}
