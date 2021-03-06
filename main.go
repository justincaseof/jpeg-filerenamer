package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/rwcarlsen/goexif/exif"
	"io"
	"jpeg-filerenamer/userconfirmation"
	"os"
	"path/filepath"
	"strings"
)

var basedir string
var successCount int64
var errorCount int64
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

	log.Info("------------------------- done.")
	log.Infof("Successfully processed %d files.", successCount)
	log.Infof("Couldn't rename %d files.", errorCount)
	log.Info("")
	log.Info("Exiting. Bye-bye")
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
			newPath = checkExisting(path, newPath)
			log.Infof("   -> new targeted file name: %s", newPath)

			// edge case: does our file already have the right name?
			if strings.Compare(path, newPath) == 0 {
				log.Info("   -> file already has the proper name!")
				return nil
			}

			// actual rename
			err = os.Rename(path, newPath)
			if err == nil {
				successCount++
			} else {
				log.ErrorF("Couldn't rename file: %s", err.Error())
				errorCount++
			}
		}
	}

	return nil
}
func checkExisting(previousfilename string, newFileNameBase string) string {
	suffixNum := 0
	result := fmt.Sprintf("%s.jpg", newFileNameBase)
	for {
		if _, err := os.Stat(result); os.IsNotExist(err) {
			break
		} else {
			if err == nil {
				// FIXME: compare file hashes here!
				fileHashesAreEqual := fileHashesAreEqual(previousfilename, result)

				if !fileHashesAreEqual {
					// increase as long as we're alreading having a file with the same props.
					suffixNum++
					log.InfoF("   -> file with same date already exists, increased suffix to %d", suffixNum)
					result = fmt.Sprintf("%s-#%d.jpg", newFileNameBase, suffixNum)
				} else {
					log.InfoF("   -> files are the same!")
					break
				}
			} else {
				log.ErrorF("Unexpected error: ", err)
				panic("Exiting...")
			}
		}
	}
	return result
}

func fileHashesAreEqual(path1 string, path2 string) bool {
	md5_1, err := md5sum(path1)
	if err != nil {
		log.Error("Could not calculate MD5 sum")
		return false
	}
	md5_2, err := md5sum(path2)
	if err != nil {
		log.Error("Could not calculate MD5 sum")
		return false
	}
	return strings.Compare(md5_1, md5_2) == 0
}

func md5sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func FormattedDateStringFromJPEGFile(fileName string) (string, error) {
	// OPEN file
	fileReader, err := os.Open(fileName)
	if err != nil {
		log.ErrorF("Error opening file", err)
		return "", err
	}
	defer fileReader.Close()

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
