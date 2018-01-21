// Cryptbin
// (C) Copyright George Saussy 2018
// version 0.1
package main

import (
  "database/sql"
  "flag"
  "fmt"
  "html/template"
  "net/http"
  _ "github.com/mattn/go-sqlite3"
)
// Structs

// ServerConfig contains the configuration of the server. 
// Should contain all of the global information for the server.
type ServerConfig struct {
  InProduction bool
  Debug bool
}

// Message struct is designed to hold a stored message. 
// The struct can also be passed to the view template to render
// the message.
type Message struct {
  Key string
  Text string
}

// Functions

// sendToLog sends information to console, or a log file if set
func sendToLog(s string, config * ServerConfig) {
  fmt.Printf("%s\n", s)
}

// logInfo prints information to console, or a log file if set
func logInfo(s string, config * ServerConfig) {
  sendToLog("INFO: " + s, config)
}

// logErrors prints error information to the log
func logError(s string, config * ServerConfig) {
  sendToLog("ERROR: " + s, config)
}

// logDebug prints debugging information to the log
func logDebug(s string, config * ServerConfig) {
  if config.Debug {
    sendToLog("DEBUG: "+ s, config)
  }
}

// checkAndInitDataBase checks if a database is present.
// If no database present, it initializes one.
func checkAndInitDataBase(config *ServerConfig) {
  if ! config.InProduction {
    db, err := sql.Open("sqlite3", "./db/cryptbin.db")
    rows, err := db.Query("PRAGMA table_info(pastes)")
    nRows int
    err := rows.Scan(&nRows)
    fmt.Printf(nRows) 
  } else {
    logError("PG not implemented yet.", config)
  }
}

// pasteHandler responds to the `/paste/` URI and sends the
// paste home page located at `html/paste.html`.
func pasteHandler(w http.ResponseWriter, r *http.Request, config * ServerConfig) {
  logInfo("Rendering paste.", config)
  t, _ := template.ParseFiles("html/paste.html")
  t.Execute(w, nil)
}

// newHandler responds to POST requests posting a new form.
// A correctly formated POST is then redirected to the view 
// the submited content.
func newHandler(w http.ResponseWriter, r *http.Request, config * ServerConfig) {
  logInfo("Computing new.", config)
  if err := r.ParseForm(); err != nil {
    logError("ParseForm() err: " + err.Error(), config)
    return
  }
  logDebug("Value: " + r.FormValue("paste_value"), config)
  logDebug("Expiration: " + r.FormValue("expiration"), config)
  http.Redirect(w, r, "/view/", 301)
}

// aboutHandler responds to the `/about/` URI and sends the
// about page for Cryptbin located at `html/about.html`.
func aboutHandler(w http.ResponseWriter, r *http.Request, config *ServerConfig) {
  logInfo("Rendering about.", config)
  t, _ := template.ParseFiles("html/about.html")
  t.Execute(w, nil)
}

// viewHandler responds to the `/view/` URI and sends the 
// view page template located at `html/view.html` using
// information given in te URI.
func viewHandler(w http.ResponseWriter, r *http.Request, config *ServerConfig) {
  logInfo("Rendering view.", config)
  dat := &Message{Key: "[tktk key]", Text: "[tktk text]"}
  t, _ := template.ParseFiles("html/view.html")
  t.Execute(w, dat)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *ServerConfig), config *ServerConfig) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fn(w, r, config)
  }
}
// main function is just the main function of this code base.
// As of right now, all that main does is 
// parse the command line arguments,
// initialize global variables, and 
// initialize an http server accordingly. 
func main() {
  // Parse flags
  inProduction :=  flag.Bool("prod", false, "Set as true if the server is running in production.")
  debug := flag.Bool("debug", false, "Set to true to print debugging information to log.")
  flag.Parse()
  config := &ServerConfig{InProduction: *inProduction, Debug: *debug}
  // Init database if need be  
  checkAndInitDataBase(config)
  // Set up http server
  http.HandleFunc("/paste/", makeHandler(pasteHandler, config)) // Route `paste` to ''
  http.HandleFunc("/new/", makeHandler(newHandler, config)) // Route `new` to ''
  http.HandleFunc("/view/", makeHandler(viewHandler, config)) // Route `view` to ''
  http.HandleFunc("/about/", makeHandler(aboutHandler, config)) // Route `about` to ''
  http.HandleFunc("/", makeHandler(pasteHandler, config)) // Route root to `paste`
  if * inProduction {
    fmt.Printf("Running Cryptbin in production mode\n")
    fmt.Printf("Go to localhost:8080\n")
    http.ListenAndServe(":8080",nil)
  } else {
    fmt.Printf("Running Cryptbin in developer mode\n")
    fmt.Printf("Go to localhost:8000\n")
    http.ListenAndServe(":8000", nil)
  }
}
