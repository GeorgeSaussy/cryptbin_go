// Cryptbin
// (C) Copyright George Saussy 2018
// version 0.1
package main

import (
  "flag"
  "fmt"
  "html/template"
  "net/http"
)

// global variables
var in_production bool

// Message struct is designed to hold a stored message. 
// The struct can also be passed to the view template to render
// the message.
type Message struct {
  Key string
  Text string
}

// pasteHandler responds to the `/paste/` URI and sends the
// paste home page located at `html/paste.html`.
func pasteHandler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("html/paste.html")
  t.Execute(w, nil)
}

// aboutHandler responds to the `/about/` URI and sends the
// about page for Cryptbin located at `html/about.html`.
func aboutHandler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("html/about.html")
  t.Execute(w, nil)
}

// viewHandler responds to the `/view/` URI and sends the 
// view page template located at `html/view.html` using
// information given in te URI.
func viewHandler(w http.ResponseWriter, r *http.Request) {
  dat := &Message{Key: "[tktk key]", Text: "[tktk text]"}
  t, _ := template.ParseFiles("html/view.html")
  t.Execute(w, dat)
}

// main function is just the main function of this code base.
// As of right now, all that main does is 
// parse the command line arguments,
// initialize global variables, and 
// initialize an http server accordingly. 
func main() {
  // Parse flags
  in_production := flag.Bool("prod", false, "set if server running in production")
  flag.Parse()
  // Set up http server
  http.HandleFunc("/paste/", pasteHandler) // Route `paste` to ''
  http.HandleFunc("/view/", viewHandler) // Route `view` to ''
  http.HandleFunc("/about/", aboutHandler) // Route `about` to ''
  http.HandleFunc("/", pasteHandler) // Route root to `paste`
  if *in_production {
    http.ListenAndServe(":8080", nil)
    fmt.Printf("Running Cryptbin in production mode\n")
    fmt.Printf("Go to localhost:8080\n")
  } else {
    http.ListenAndServe(":8000", nil)
    fmt.Printf("Running Cryptbin in developer mode\n")
    fmt.Printf("Go to localhost:8000\n")
  }
}
