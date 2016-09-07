package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Event encapsulates the details of a person
type Event struct {
	Name      string
	Sex       string
	Age       int
	Timestamp int64
}

var (
	file *os.File
	err  error
	Log  *log.Logger
)

func check(e error) {
	if e != nil {
		Log.Printf("Received unrecoverable error, %v", e)
		panic(e)
	}
}

func setup() {
	dir := "/var/segment/app/contents"
	_, err := os.Stat(dir)
	if err != nil {
		// Create the content dir
		os.Mkdir(dir, 0744)
		Log.Println("Creating the content directory")
	}
}

func streamWriter() {
	// Opening contents file for writing
	file, err = os.OpenFile("/var/segment/app/contents/contents.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Log.Println("Could not open contents file")
		panic(err)
	}

	defer file.Close()

	files := [1]string{"/var/segment/app/seed/yob2015.txt"}
	for _, f := range files {
		Log.Printf("Reading new file %s\n", f)

		// Open each file in the seed directory for reading
		rfile, err := os.Open(f)
		check(err)
		defer rfile.Close()

		// Open the scanner to read from the file
		scanner := bufio.NewScanner(rfile)
		var buf bytes.Buffer

		// Continue using the contents from the file till we reach EOF
		for scanner.Scan() {
			text := scanner.Text() // Each line extracted from the seed file
			tokens := strings.Split(text, ",")
			age, _ := strconv.Atoi(tokens[2])
			event := Event{tokens[0], tokens[1], age, time.Now().UnixNano()}
			b, err := json.Marshal(event)
			if err == nil {
				buf.Write(b)
				buf.Write([]byte("\n"))

				file.Write(buf.Bytes())
				time.Sleep(2 * time.Millisecond)
				buf.Reset()
			}
		}
	}
	Log.Println("Completed writing files")
}

func main() {

	// Opening log file for writing
	f, err := os.OpenFile("/var/segment/log/streaming-app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Log.Println("Could not open log file")
		panic(err)
	}
	defer f.Close()

	Log = log.New(f, "INFO", log.Ldate|log.Ltime|log.Lshortfile)

	setup()

	go streamWriter()

	Log.Println("Application started")

	// This is just to keep a long running app
	http.ListenAndServe(":8080", nil)
}
