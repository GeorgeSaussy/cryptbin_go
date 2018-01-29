package main

import (
 "hash/fnv"
 "strconv"
 "strings"
 "time"
 "errors"
)

// checkAndInitDataBase checks if a database is present.
// If no database present, it initializes one.
func checkAndInitDataBase(config *ServerConfig) {
  if config.Database != nil {
    logDebug("Database initialized!", config)
  }
}

// hashString wraps the FNV hash function to hash strings
func hashString(s string) string {
  ret := ""
  h := fnv.New64a()
  h.Write([]byte(s))
  ret = strconv.FormatUint(h.Sum64(), 16)
  return ret
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

// getExpirationTime get expiration time from the POST parameter
func getExpirationTime(expiration_text string) time.Time {
  ret := time.Now()
  switch expiration_text {
  case "burn_after_reading":
    ret = time.Now().Add(time.Hour*12)
  case "12_hr":
    ret = time.Now().Add(time.Hour*12)
  case "24_hr":
    ret = time.Now().Add(time.Hour*24)
  case "5_days":
    ret = time.Now().Add(time.Hour*5*24)
  default:
    ret = time.Now().Add(time.Hour*12) // Default to 12 hr 
  }
  return ret
}

// storeNewPaste stores a paste in the database 
func storeNewPaste(created_at time.Time, paste_value string, expiration time.Time, burn_after_reading bool, paste_hash string, config *ServerConfig) {
  stmt, _ := config.Database.Prepare("INSERT INTO paste(created_at, paste_value, expiration, burn_after_reading, paste_hash, view_count) values(?, ?, ?, ?, ?, ?)")
  stmt.Exec(created_at, paste_value, expiration, burn_after_reading, paste_hash, 0)
}

// queryForPasteValue runs a query on the database to get the paste value associated by the hash
func queryForPasteValue(paste_hash string, config *ServerConfig) (string, error) {
  ret := ""
  rows, err := config.Database.Query("SELECT paste_value FROM paste WHERE paste_hash='" + paste_hash + "'")
  for rows.Next() {
    rows.Scan(&ret)
  }
  rows.Close()
  config.Database.Exec("UPDATE paste SET view_count = view_count + 1 WHERE paste_hash='" + paste_hash + "'")
  if ret == "" || err != nil {
    err = errors.New("Hash not in database.")
  }
  return ret, err
}

// cleanDatabase deletes all rows from the database that are too old
func cleanDatabase(config *ServerConfig) {
  config.Database.Exec("DELETE FROM paste WHERE expiration < DATETIME('now')")
  config.Database.Exec("DELETE FROM paste WHERE view_count > 1 AND burn_after_reading = 1")
}

func checkIfBurnAfterReading(paste_hash string, config *ServerConfig) bool {
  ret := false
  rows, _ := config.Database.Query("SELECT burn_after_reading FROM paste WHERE paste_hash='" + paste_hash + "'")
  for rows.Next() {
    rows.Scan(&ret)
  }
  rows.Close()
  return ret
}

func getExpiration(paste_hash string, config *ServerConfig) time.Time {
  ret := time.Now()
  rows, _ := config.Database.Query("SELECT expiration FROM paste WHERE paste_hash='" + paste_hash + "'")
  for rows.Next() {
    rows.Scan(&ret)
  }
  rows.Close()
  return ret
}

// getTimeLeftMsg get the message to alert users how long a message
// has left to live
func getTimeLeftMsg(paste_hash string, config *ServerConfig) string {
  ret := ""
  burn_after_reading := checkIfBurnAfterReading(paste_hash, config)
  expiration := getExpiration(paste_hash, config)
  if ! burn_after_reading {
    ret = strconv.FormatFloat(expiration.Sub(time.Now()).Hours(), 'f', 2, 64)
    ret = ret + " hours left to read paste."
  } else {
    if burn_after_reading {
      ret = "Paste set to burn after reading."
    }
  }
  return ret
}



