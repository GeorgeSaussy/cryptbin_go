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
  "time"
  _ "github.com/lib/pq"
  _ "github.com/mattn/go-sqlite3"
)

// Structs

// ServerConfig contains the configuration of the server. 
// Should contain all of the global information for the server.
type ServerConfig struct {
  InProduction bool
  Debug bool
  Database *sql.DB
}

// Message struct is designed to hold a stored message. 
// The struct can also be passed to the view template to render
// the message.
type Message struct {
  Key string
  Text string
  TimeLeftMsg string
}

// Global variables
var config *ServerConfig

// Functions

// getBlankConfig returns a pointer to a blank ServerConfig object
func getBlankConfig() *ServerConfig {
  ret :=&ServerConfig{InProduction: false, Debug: false, Database: nil}
  return ret
}

// pasteHandler responds to the `/paste/` URI and sends the
// paste home page located at `html/paste.html`.
func pasteHandler(w http.ResponseWriter, r *http.Request) {
  logInfo("Rendering paste.")
  t, _ := template.ParseFiles("html/paste.html")
  t.Execute(w, nil)
}

// newHandler responds to POST requests posting a new form.
// A correctly formated POST is then redirected to the view 
// the submited content.
func newHandler(w http.ResponseWriter, r *http.Request) {
  logInfo("Computing new.")
  created_at := time.Now()
  paste_value := r.FormValue("paste_value")
  expiration_text := r.FormValue("expiration")
  expiration := getExpirationTime(expiration_text)
  burn_after_reading := (expiration_text == "burn_after_reading")
  paste_hash := hashString(paste_value)
  paste_key := r.FormValue("paste_key")
  logDebug("Value: "+paste_value+"\nExpiration: "+expiration_text+"\nHash: "+paste_hash+"\nKey: "+paste_key)
  storeNewPaste(created_at, paste_value, expiration, burn_after_reading, paste_hash)
  http.Redirect(w, r, "/view/"+paste_hash+"#"+paste_key, 301)
}

// aboutHandler responds to the `/about/` URI and sends the
// about page for Cryptbin located at `html/about.html`.
func aboutHandler(w http.ResponseWriter, r *http.Request) {
  logInfo("Rendering about.")
  t, _ := template.ParseFiles("html/about.html")
  t.Execute(w, nil)
}

// viewHandler responds to the `/view/` URI and sends the 
// view page template located at `html/view.html` using
// information given in te URI.
func viewHandler(w http.ResponseWriter, r *http.Request) {
  logInfo("Rendering view.")
  paste_hash := r.URL.Path[len("/view/"):]
  if checkValidHash(paste_hash) {
    paste_value, err := queryForPasteValue(paste_hash)
    time_left_msg := getTimeLeftMsg(paste_hash)
    if err != nil {
      _404Handler(w, r)
    } else {
      logDebug("hash: "+paste_hash+"\npaste: "+paste_value)
      cleanDatabase()
      dat := &Message{Key: paste_hash, Text: paste_value, TimeLeftMsg: time_left_msg}
      t, _ := template.ParseFiles("html/view.html")
      t.Execute(w, dat)
    }
  } else {
    _404Handler(w, r)
  }
}

// _404Handler routes the user to a 404 page
func _404Handler(w http.ResponseWriter, r *http.Request) {
  logInfo("Rendering 404.")
  t, _ := template.ParseFiles("html/404.html")
  t.Execute(w, nil)
}

// makeHandler wraps the HandlerFuncs to pass the server configuration.
func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fn(w, r)
  }
}

// main function is just the main function of this code base.
// As of right now, all that main does is 
// parse the command line arguments,
// initialize global variables, and 
// initialize an http server accordingly. 
func main() {
  // Parse flags
  in_production :=  flag.Bool("prod", false, "Set as true if the server is running in production.")
  debug := flag.Bool("debug", false, "Set to true to print debugging information to log.")
  db_password := flag.String("pwd", "", "Set the Postgres database password")
  db_username := flag.String("uname", "", "Set the Postgres username.")
  db_name := flag.String("dbname", "", "Set the Postgres database name.")
  flag.Parse()
  config = &ServerConfig{InProduction: *in_production, Debug: *debug, Database: nil}
  // Init database if need be
  if * in_production {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", *db_username, *db_password, *db_name)
    fmt.Printf(dbinfo)
    db_instance, err := sql.Open("postgres", dbinfo)
    if err != nil {
      logError(err.Error())
    }
    config.Database = db_instance
  } else {
    config.Database, _ = sql.Open("sqlite3", "db/cryptbin.db")
  }
  // Set up http server
  http.HandleFunc("/paste/", makeHandler(pasteHandler)) // Route `paste` to ''
  http.HandleFunc("/new/", makeHandler(newHandler)) // Route `new` to ''
  http.HandleFunc("/view/", makeHandler(viewHandler)) // Route `view` to ''
  http.HandleFunc("/about/", makeHandler(aboutHandler)) // Route `about` to ''
  http.HandleFunc("/", makeHandler(pasteHandler)) // Route root to `paste`
  http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("html/assets")))) // Route `assets/*` to `html/assets` 
  if * in_production {
    fmt.Printf("Running Cryptbin in production mode\n")
    fmt.Printf("Go to localhost:8080\n")
    err := http.ListenAndServe(":8080", nil)
    fmt.Printf(err.Error())
    fmt.Printf("\n")
  } else {
    fmt.Printf("Running Cryptbin in developer mode\n")
    fmt.Printf("Go to localhost:8000\n")
    err := http.ListenAndServe(":8000", nil)
    fmt.Printf(err.Error())
    fmt.Printf("\n")
  }
}
