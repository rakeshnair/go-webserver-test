package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Page struct {
	Title string
	Body  []byte
}

// Event encapsulates the details of a person
type Event struct {
	Name      string
	Sex       string
	Age       int
	Timestamp int64
}

const (
	contentDir  string = "/var/segment/app/contents/"
	seedDir     string = "/var/segment/app/seed/"
	contentFile string = "contents.log"
	logFile     string = "/var/segment/log/main.app.log"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	WarnLogger  *log.Logger
)

func initialize() {
	_, err := os.Stat(contentDir)
	if err != nil {
		// Create the content dir
		os.Mkdir(contentDir, 0744)
		InfoLogger.Println("Creating the content directory")
	}
	_, err = os.Stat(contentDir + contentFile)
	if err != nil {
		InfoLogger.Println("Creating the content file")
		os.Create(contentDir + contentFile)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (p *Page) save() {
	filename := contentDir + contentFile
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	if err != nil {
		ErrorLogger.Println("Could not open file to write")
	} else {
		_, err := f.Write(p.Body)
		if err != nil {
			ErrorLogger.Printf("Could not write to file. Error: %s", err.Error())
		}
	}
	InfoLogger.Println("Contents updated")
}

func loadPage() (*Page, error) {
	filename := contentDir + contentFile
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{"contents", body}, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	fmt.Fprintf(w, "Hello %s", name)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	page, err := loadPage()
	if err == nil {
		fmt.Fprintf(w, "%s", page.Body)
		InfoLogger.Println("Contents displayed on screen")
	} else {
		fmt.Fprintf(w, "Page not found")
		ErrorLogger.Println("Page not found. Contents could not be displayed")
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(contentDir)
	if err != nil {
		fmt.Fprintf(w, "Contents cannot be listed")
		ErrorLogger.Println("Contents cannot be listed on the webpage")
	} else {
		for _, e := range files {
			fmt.Fprintf(w, "%s", e.Name())
		}
	}
	InfoLogger.Println("List of content files are displayed on the webpage")
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	randomFile := randomFile(seedDir)
	pageBytes := seedContentToBytes(seedDir + randomFile)
	page := &Page{randomFile[:strings.LastIndex(randomFile, ".")], pageBytes}
	page.save()

	fmt.Fprintf(w, "Page %s saved", randomFile[:strings.LastIndex(randomFile, ".")])

}

func seedContentToBytes(input string) []byte {
	rfile, err := os.Open(input)
	check(err)
	defer rfile.Close()

	scanner := bufio.NewScanner(rfile)
	var buf bytes.Buffer

	for scanner.Scan() {
		text := scanner.Text()
		tokens := strings.Split(text, ",")
		age, _ := strconv.Atoi(tokens[2])
		event := Event{tokens[0], tokens[1], age, time.Now().UnixNano()}
		b, err := json.Marshal(event)
		if err == nil {
			buf.Write(b)
			buf.Write([]byte("\n"))
		}
	}

	return buf.Bytes()
}

func randomFile(dirname string) string {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	files, _ := ioutil.ReadDir(dirname)
	randomIndex := random.Intn(len(files))
	return files[randomIndex].Name()
}

func sigHUPHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for sig := range c {
			WarnLogger.Printf("Got a HUP signal(%v). Ignoring for now", sig)
		}
	}()
}

func main() {
	// Opening log file for writing
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Could not open log file")
		panic(err)
	}
	defer f.Close()

	InfoLogger = log.New(f, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(f, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(f, "INFO", log.Ldate|log.Ltime|log.Lshortfile)

	initialize()

	InfoLogger.Println("Application started")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/write", writeHandler)

	sigHUPHandler()

	http.ListenAndServe(":8080", nil)
}
