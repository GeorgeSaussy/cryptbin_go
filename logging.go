package main

import "fmt"

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
