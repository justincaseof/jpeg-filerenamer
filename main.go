package main

import (
	"flag"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/rwcarlsen/goexif/exif"
	"jpeg-filerenamer/userconfirmation"
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

	log.Info("#########################")
	log.Info("### JPEG file renamer ###")
	log.Info("#########################")
	log.Info("")

	flag.StringVar(&basedir, "source", "", "The directory to traverse for JPEG files.")
	flag.Parse()

	log.ErrorF("WE'RE ABOUT TO RENAME *ALL* JPG/JPEG FILES IN GIVEN DIRECTORY '%s'...", basedir)
	log.ErrorF(">>>> ARE YOU ***ABSOLUTELY*** SURE? <<<<")
	if !userconfirmation.AskForConfirmation() {
		log.Info("Aborted operation. Bye-bye")
		return
	}

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
			dir := strings.TrimSuffix(path, info.Name())
			newPath := dir + dateString
			// prepend a suffix in case a file with the same name already exists
			newPath = checkExisting(newPath)
			log.Infof("   -> new targeted file name: %s", newPath)
			count++
		}

		// TODO: rename file
	}

	return nil
}
func checkExisting(newFileNameBase string) string {
	suffixNum := 0
	result := fmt.Sprintf("%s.jpg", newFileNameBase)
	for {
		if _, err := os.Stat(result); os.IsNotExist(err) {
			break
		} else {
			if err == nil {
				// increase as long as we're alreading having a file with the same props.
				suffixNum++
				log.InfoF("   -> file with same date already exists, increased suffix to %d", suffixNum)
				result = fmt.Sprintf("%s-#%d.jpg", newFileNameBase, suffixNum)
			} else {
				log.ErrorF("Unexpected error: ", err)
				panic("Exiting...")
			}
		}
	}
	return result
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

	return time.Format("2006-01-02_15-04-05"), nil
}
