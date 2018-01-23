// Cryptbin
// (C) Copyright George Saussy 2018
// version 0.1
package main

import (
  "database/sql"
  "flag"
  "fmt"
  "hash/fnv"
  "html/template"
  "net/http"
  "strconv"
  "strings"
  "time"
  _ "github.com/mattn/go-sqlite3"
)

// Structs

// ServerConfig contains the configuration of the server. 
// Should contain all of the global information for the server.
type ServerConfig struct {
  InProduction bool
  Debug bool
  Password string
  Username string
  Database *sql.DB
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

func checkError(err error) {
  if err != nil {
    panic(err)
  }
}

// checkAndInitDataBase checks if a database is present.
// If no database present, it initializes one.
func checkAndInitDataBase(config *ServerConfig) {
  if config.Database != nil {
    logDebug("Database initialized!", config)
  }
}

// pasteHandler responds to the `/paste/` URI and sends the
// paste home page located at `html/paste.html`.
func pasteHandler(w http.ResponseWriter, r *http.Request, config * ServerConfig) {
  logInfo("Rendering paste.", config)
  t, _ := template.ParseFiles("html/paste.html")
  t.Execute(w, nil)
}

// hashString wraps the FNV hash function to hash strings
func hashString(s string) string {
  ret := ""
  h := fnv.New64a()
  h.Write([]byte(s))
  ret = strconv.FormatUint(h.Sum64(), 16) // XXX check conversion for correctness
  return ret
}

// newHandler responds to POST requests posting a new form.
// A correctly formated POST is then redirected to the view 
// the submited content.
func newHandler(w http.ResponseWriter, r *http.Request, config * ServerConfig) {
  logInfo("Computing new.", config)
  created_at := time.Now()
  paste_value := r.FormValue("paste_value")
  expiration_text := r.FormValue("expiration")
  expiration := time.Now()
  burn_after_reading := false
  switch expiration_text {
  case "burn_after_reading":
    burn_after_reading = true
    expiration = expiration.Add(time.Hour*12)
  case "12_hr":
    expiration = expiration.Add(time.Hour*12)
  case "24_hr":
    expiration = expiration.Add(time.Hour*24)
  case "5_days":
    expiration = expiration.Add(time.Hour*5*24)
  default:
    expiration = expiration.Add(time.Hour*12) // Default to 12 hr 
  }
  paste_hash := hashString(paste_value)
  logDebug("Value: " + paste_value, config)
  logDebug("Expiration: " + expiration_text, config)
  logDebug("Hash: " + paste_hash, config)
  stmt, _ := config.Database.Prepare("INSERT INTO paste(created_at, paste_value, expiration, burn_after_reading, paste_hash) values(?, ?, ?, ?, ?)")
  stmt.Exec(created_at, paste_value, expiration, burn_after_reading, paste_hash)
  http.Redirect(w, r, "/view/"+paste_hash, 301)
}

// aboutHandler responds to the `/about/` URI and sends the
// about page for Cryptbin located at `html/about.html`.
func aboutHandler(w http.ResponseWriter, r *http.Request, config *ServerConfig) {
  logInfo("Rendering about.", config)
  t, _ := template.ParseFiles("html/about.html")
  t.Execute(w, nil)
}

// checkValidHash checks if a string is a hexidecimal string
func checkValidHash(paste_hash string) bool {
  ret := true
  valid_chars := "0123456789abcdefABCDEF"
  for ind := 0; ind < len(paste_hash) && ret; ind++ {
    if ! strings.Contains(valid_chars, paste_hash[ind:ind+1]) {
      ret = false
    }
  }
  return ret
}

// viewHandler responds to the `/view/` URI and sends the 
// view page template located at `html/view.html` using
// information given in te URI.
func viewHandler(w http.ResponseWriter, r *http.Request, config *ServerConfig) {
  logInfo("Rendering view.", config)
  uri := r.URL.Path
  paste_hash := uri[len("/view/"):]
  if checkValidHash(paste_hash) {
    rows, _ := config.Database.Query("SELECT paste_value FROM paste WHERE paste_hash='" + paste_hash + "'")
    paste_value := "[tktk]"
    for rows.Next() {
      rows.Scan(&paste_value)
    }
    rows.Close()
    logDebug("URI: " + uri, config)
    logDebug("hash: " + paste_hash, config)
    logDebug("paste: " + paste_value, config)
    dat := &Message{Key: paste_hash, Text: paste_value}
    t, _ := template.ParseFiles("html/view.html")
    t.Execute(w, dat)
  } else { // in case hash is invalid, redirect home
    // FIXME implement a 404 page
    pasteHandler(w, r, config)
  }
}

// makeHandler wraps the HandlerFuncs to pass the server configuration.
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
  password := flag.String("pwd", "", "Set the Postgres database password")
  username := flag.String("uname", "", "Set the Postgres username.")
  flag.Parse()
  config := &ServerConfig{InProduction: *inProduction, Debug: *debug, Password: *password, Username: *username, Database: nil}
  // Init database if need be  
  db, err := sql.Open("sqlite3", "db/cryptbin.db")
  if err != nil {
    logError("Database is busted", config)
  } else {
    config.Database = db
  }
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
