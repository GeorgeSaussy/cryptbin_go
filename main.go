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
  created_at := time.Now()
  paste_value := r.FormValue("paste_value")
  expiration_text := r.FormValue("expiration")
  expiration := getExpirationTime(expiration_text)
  burn_after_reading := (expiration_text == "burn_after_reading")
  paste_hash := hashString(paste_value)
  logDebug("Value: "+paste_value+"\nExpiration: "+expiration_text+"\nHash: "+paste_hash, config)
  storeNewPaste(created_at, paste_value, expiration, burn_after_reading, paste_hash, config)
  http.Redirect(w, r, "/view/"+paste_hash, 301)
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
  paste_hash := r.URL.Path[len("/view/"):]
  if checkValidHash(paste_hash) {
    paste_value, err := queryForPasteValue(paste_hash, config)
    if err != nil {
      _404Handler(w, r, config)
    } else {
      logDebug("hash: "+paste_hash+"\npaste: "+paste_value, config)
      cleanDatabase(config)
      dat := &Message{Key: paste_hash, Text: paste_value}
      t, _ := template.ParseFiles("html/view.html")
      t.Execute(w, dat)
    }
  } else {
    _404Handler(w, r, config)
  }
}

// _404Handler routes the user to a 404 page
func _404Handler(w http.ResponseWriter, r *http.Request, config *ServerConfig) {
  logInfo("Rendering 404.", config)
  t, _ := template.ParseFiles("html/404.html")
  t.Execute(w, nil)
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
