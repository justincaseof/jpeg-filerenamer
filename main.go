package main

import (
	"flag"
	"github.com/apsdehal/go-logger"
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"path/filepath"
	"strings"
)

var basedir string
var count int64
var log *logger.Logger

func main() {
	// Init Logger
	var err error
	log, err = logger.New("test", 1, os.Stdout)
	if err != nil {
		panic(err)
	}

	log.Info("### JPEG file renamer")

	flag.StringVar(&basedir, "source", "", "The directory to traverse for JPEG files.")
	flag.Parse()

	log.Info("Args:")
	log.InfoF("  -> basedir: %s", basedir)

	// some validation and the final check for a required succeeding "/".
	if len(basedir) < 1 {
		log.Error("no path given!")
		return
	}
	if _, err := os.Stat(basedir); os.IsNotExist(err) {
		log.Error("given path does not exist!")
		return
	}
	if !strings.HasSuffix(basedir, "/") {
		basedir = basedir + "/"
	}

	filepath.Walk(basedir, checkJPEG)

	log.Infof("------------------------- done.", count)
	log.Infof("Processed %d files.", count)
}

func checkJPEG(path string, info os.FileInfo, err error) error {
	if strings.HasSuffix(info.Name(), ".jpg") || strings.HasSuffix(info.Name(), ".jpeg") {
		log.Infof("* Checking file %s", path)
		// found a potential candidate! let's throw it upon our exif parser
		dateString, err := FormattedDateStringFromJPEGFile(path)
		if err != nil {
			log.Warning("  -> could not get date from JPEG file!")
		} else {
			log.Infof("  -> Found date: %s", dateString)
			count++
		}

		// TODO: rename file
	}

	return nil
}

func FormattedDateStringFromJPEGFile(fileName string) (string, error) {
	// OPEN file
	fileReader, err := os.Open(fileName)
	if err != nil {
		log.ErrorF("Error opening file", err)
		return "", err
	}

	// DECODE file
	exif, err := exif.Decode(fileReader)
	if err != nil {
		log.ErrorF("Error decoding file", err)
		return "", err
	}

	// GET date
	time, err := exif.DateTime()
	if err != nil {
		log.ErrorF("Error retrieving date from file", err)
		return "", err
	}
	log.InfoF("  -> found time: %v", time)

	//return time.Format("yyyy-MM-dd_HH-mm-ss"), nil
	return time.Format("2006-01-02_15-04-05"), nil
}
