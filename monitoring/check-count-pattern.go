package main

/*
	Count pattern present in a file
	and raise alert base on warning value
	and critical value.
*/

import (
	"bufio"
	"fmt"
	"flag"
	"os"
	"path"
	"io/ioutil"
	"strconv"
	"regexp"
)

func checkError(e error) {
	/*
		Create error and return
		fatal message is error present
	*/
	if e != nil {
		panic(e)
	}
}

func man() {
	/* 
		Create man output.
	*/
	fmt.Printf("\n%s -logfile=/Path/to/filename.log -offset=/Path/to/offset/file -pattern=\"regexp[a-z]\" -warn=2 -crit=8", path.Base(os.Args[0]))
	fmt.Println("\n")
	flag.PrintDefaults()
	fmt.Println("\n")

	// Exit with error statu
	os.Exit(1)
}


func main() {
	// Define Variable
	matchCounter := 0

	// Create arguments parser.
	logFileName := flag.String("logfile", "", "Log file name.")
	offsetFileName := flag.String("offset", "", "Offset file name.")
	logPattern := flag.String("pattern", "", "Define matching pattern.")
	warningLimit := flag.Int("warn", 0, "Define warning limit.")
	criticalLimit := flag.Int("crit", 0, "Define critical limit.")

	flag.Parse()

	// Check if required arguments is define.
	if *logFileName == "" || *logPattern == "" {
		man()
	}

	// Try to open file.
	fileHandle, err := os.Open(*logFileName)
	checkError(err)
	defer fileHandle.Close()

	// Get current file size
	fileStat, err := fileHandle.Stat()
	checkError(err)

	/*
		If offset file is define, get offset value from last run
		and start to read from last offset.
	*/
	if *offsetFileName != "" {

		// If offset file exist.
		if _, err := os.Stat(*offsetFileName); ! os.IsNotExist(err) {

			// Get offset from offset file.
			offset, err := ioutil.ReadFile(*offsetFileName)
			checkError(err)

			// Set seek cursor value.
			seek_cursor, _ := strconv.ParseInt(string(offset), 10, 64)

			// If current size is bigger than offset
			if fileStat.Size() > seek_cursor {
				_, err = fileHandle.Seek(seek_cursor, 0)
				checkError(err)
			} 
		} 


	}

	fileScanner := bufio.NewScanner(fileHandle)

	// Compile the expression
	regPattern := regexp.MustCompile(*logPattern)

	// for each line
	for fileScanner.Scan() {
		if regPattern.MatchString(fileScanner.Text()) {
			matchCounter++
		}
	}

	// Save offset to offset file.
	offset := []byte(strconv.Itoa(int(fileStat.Size())))
	err = ioutil.WriteFile(*offsetFileName, offset, 0644)
	checkError(err)
	
	// Exit base on warn and criti limit
	if matchCounter >= *criticalLimit {
		// Critical limit reach
		fmt.Printf("[CRITICAL] - %d match found in %s. (%d > %d)", matchCounter, *logFileName, matchCounter, *criticalLimit)
		os.Exit(2)
	} else if *criticalLimit > matchCounter && matchCounter >= *warningLimit {
		// Match count between critical and warn limit.
		fmt.Printf("[WARNING] - %d match found in %s. (%d > %d)", matchCounter, *logFileName, matchCounter, *warningLimit)
		os.Exit(1)
	} else {
		// No alert
		fmt.Printf("[OK] - %d match found in %s. (%d < %d)", matchCounter, *logFileName, matchCounter, *warningLimit)

	}

}
