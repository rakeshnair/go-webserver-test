package main

import (
	"bufio"
	"encoding/json"
	"gopkg.in/natefinch/lumberjack.v2"
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
	Log      *log.Logger
	EventLog *lumberjack.Logger
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
	EventLog = &lumberjack.Logger{
		Filename:   "/var/segment/app/contents/lumberjack-contents.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}
}

func write(msg []byte) {
	EventLog.Write(msg)
	EventLog.Write([]byte("\n"))

}

func streamWriter() {
	files := [1]string{"/var/segment/app/seed/yob2015.txt"}
	for _, f := range files {
		Log.Printf("Reading new file %s\n", f)

		// Open each file in the seed directory for reading
		rfile, err := os.Open(f)
		check(err)
		defer rfile.Close()

		// Open the scanner to read from the file
		scanner := bufio.NewScanner(rfile)

		// Continue using the contents from the file till we reach EOF
		for scanner.Scan() {
			text := scanner.Text() // Each line extracted from the seed file
			tokens := strings.Split(text, ",")
			age, _ := strconv.Atoi(tokens[2])
			event := Event{tokens[0], tokens[1], age, time.Now().UnixNano()}
			b, err := json.Marshal(event)
			if err == nil {
				write(b)
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
	Log.Println("Completed writing files")
}

func main() {

	// Opening log file for writing
	lf, e1 := os.OpenFile("/var/segment/log/streaming-app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if e1 != nil {
		Log.Println("Could not open log file")
		panic(e1)
	}
	defer lf.Close()

	Log = log.New(lf, "INFO", log.Ldate|log.Ltime|log.Lshortfile)

	setup()

	go streamWriter()

	Log.Println("Application started")

	// This is just to keep a long running app
	http.ListenAndServe(":8080", nil)
}
