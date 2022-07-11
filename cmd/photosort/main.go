package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/rwcarlsen/goexif/exif"
)

var (
	sourceFolder      *string
	destinationFolder *string
	totalSize         int64 = 0
)

func main() {
	// TODO: Test if using multiple CPUs actually improves performance since  this is a sequential process
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create new parser object
	parser := argparse.NewParser("photosort", "Sorts photos from one directory to another")
	// Create string flag
	sourceFolder = parser.String("s", "source-folder", &argparse.Options{Required: true, Help: "Source folder"})
	destinationFolder = parser.String("d", "destination-folder", &argparse.Options{Required: true, Help: "Destination folder with archived sorted photos"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}
	if *sourceFolder == "" {
		log.Fatal("source must be specified")
	}
	if *destinationFolder == "" {
		log.Fatal("archive must be specified")
	}

	err = filepath.Walk(*sourceFolder, walkFunc)
	if err != nil {
		log.Println(err)
	}
}

func walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	copiedBytes, err := processFile(path, *destinationFolder)
	if err != nil {
		// we won't return nil just to not quit the walk function
		log.Println(err)
	} else {
		totalSize += copiedBytes
		log.Println(ByteCountSI(totalSize))
	}

	return nil
}

//  Converts a size in bytes to a human-readable string in SI (decimal)
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func processFile(filePath string, archiveFolder string) (int64, error) {
	filename := filepath.Base(filePath)

	date, err := getDate(filePath)
	if err != nil {
		if err := createDir(fmt.Sprintf("%s/%s", archiveFolder, "others")); err != nil {
			return 0, err
		}
		extension := filepath.Ext(filePath)
		if isImageOrVideo(extension) {
			finalPath := fmt.Sprintf("%s/%s/%s", archiveFolder, "others", filename)
			return copyFile(filePath, finalPath)
		}
		return 0, &ErrNotMediaFile{extension, filePath}
	}

	finalPath, err := newPath(archiveFolder, filename, date)
	if err != nil {
		return 0, err
	}

	return copyFile(filePath, finalPath)
}

func isImageOrVideo(extension string) bool {

	imageExtensions := map[string]bool{".tiff": true, ".tif": true, ".gif": true,
		".jpeg": true, ".jpg": true, ".png": true,
		".raw": true, ".dng": true,
		".webm": true, ".mkv": true, ".avi": true, ".mov": true, ".wmv": true,
		".mp4": true, ".MP4": true, ".m4v": true, ".mpg": true, ".mp2": true, ".mpeg": true}

	return imageExtensions[strings.ToLower(extension)]
}

func getDate(filepath string) (time.Time, error) {
	var dt time.Time
	file, err := os.Open(filepath)
	if err != nil {
		return dt, err
	}

	data, err := exif.Decode(file)
	if err != nil {
		return dt, err
	}

	return data.DateTime()
}

// Returns true if a dir/file already exists
func Exists(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// Generates the entire new path based on all the data, checks for collisions (and rename if needed)
func newPath(archive string, oldName string, date time.Time) (string, error) {
	dir := fmt.Sprintf("%s/%0004d/%0004d-%02d/%0004d-%02d-%02d", archive, date.Year(), date.Year(), date.Month(),
		date.Year(), date.Month(), date.Day())
	if err := createDir(dir); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", dir, oldName), nil
}

// Creates a directory if it doesn't exist yet
func createDir(dir string) error {
	if exists, err := Exists(dir); err != nil {
		return err
	} else if !exists {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) (int64, error) {

	if exists, err := Exists(dst); err != nil {
		return 0, err
	} else if exists {
		return 0, &ErrFileExists{filePath: dst}
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, &ErrFileNotRegular{filePath: src}
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func() { _ = source.Close() }()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() { _ = destination.Close() }()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
