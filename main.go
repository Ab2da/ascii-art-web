package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

var templatePath = "template/index.html"

func main() {
	Start()
}

// to populate the font value(banner) and the text (arguement) value for ascii-art(to choose which file)
func ascii_finder(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		showError(w, "400 BAD REQUEST", http.StatusBadRequest)
		return
	}

	font := r.FormValue("font")
	text := r.FormValue("request")
	font = font + ".txt"
	if text == "" || font == "" || strings.Contains(text, "£") {
		showError(w, "404 BANNER NOT FOUND", http.StatusNotFound)
		return
	}
	args := strings.Split(text, "\r\n")
	// ascii-art
	for _, word := range args {
		for i := 0; i < 8; i++ {

			for _, letter := range word {
				returna := GetLine((1 + int(letter-' ')*9 + i), font, w)
				if returna == "abort" {
					//not an internal error!! fix this - bad request - **
					showError(w, "400 BAD REQUEST", http.StatusBadRequest)
					return
				} else {
					// if letter < 32 || letter > ? showError(w, "400 BAD REQUEST", http.StatusBadRequest)
					// else {}
					// prints contents to web page
					fmt.Fprint(w, returna)
				}

			}
			fmt.Fprintln(w)
		}
	}
}

// similar to readfile
func GetLine(num int, filename string, w http.ResponseWriter) string {
	f, e := os.Open(filename)
	if e != nil {
		return "abort"
	}
	scanner := bufio.NewScanner(f)
	lineNum := 0
	line := ""
	for scanner.Scan() {
		if lineNum == num {
			line = scanner.Text()
		}
		lineNum++
	}
	return line
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && r.URL.Path != "/" {
		showError(w, "400 BAD REQUEST", http.StatusBadRequest)
		// return here will stop executing this function
		return
	}
	// Render  index.html template - server not working(difficult to test)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		showError(w, "404 TEMPLATE NOT FOUND", http.StatusNotFound)
		return
	}

	err = t.Execute(w, nil)

	if err != nil {
		showError(w, "500 INTERNAL SERVER ERROR4", http.StatusInternalServerError)
		return
	}
}

// Render the error.html template - when there is an error function is run
func showError(w http.ResponseWriter, message string, statusCode int) {
	t, err := template.ParseFiles("template/error.html")
	if err == nil {
		w.WriteHeader(statusCode)
		t.Execute(w, message)
	}
}

// func remove(slice []string, s int) []string {
// 	return append(slice[:s], slice[s+1:]...)
// }

// this starts the server
func Start() {
	LocalHost := "1304"
	backgroundImage := http.FileServer(http.Dir("./background"))
	exec.Command("open", "http://localhost:"+LocalHost).Start() //opens webpage
	http.Handle("/background/", http.StripPrefix("/background", backgroundImage))
	http.HandleFunc("/", formHandler) // aa way of running formhandler/ascii finder
	http.HandleFunc("/ascii-art", ascii_finder)
	fmt.Println("Server started on Port", LocalHost)

	err := http.ListenAndServe(":"+LocalHost, nil)
	if err != nil {
		log.Fatal(err)
	}
}
