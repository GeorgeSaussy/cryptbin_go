package main

import "fmt"

// sendToLog sends information to console, or a log file if set
func sendToLog(s string) {
  fmt.Printf("%s\n", s)
}

// logInfo prints information to console, or a log file if set
func logInfo(s string) {
  sendToLog("INFO: " + s)
}

// logErrors prints error information to the log
func logError(s string) {
  sendToLog("ERROR: " + s)
}

// logDebug prints debugging information to the log
func logDebug(s string) {
  if config.Debug {
    sendToLog("DEBUG: "+ s)
  }
}

// checkError checks if an error is not nil and panics
// if there is a problem
func checkError(err error) {
  if err != nil {
    panic(err)
  }
}
