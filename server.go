package main

import (
  "fmt"
  "flag"
)

var in_production bool

func main() {
  in_production := flag.Bool("prod", false, "set if server running in production")
  flag.Parse()

  if *in_production {
    fmt.Printf("Running Cryptbin in production mode\n")
  } else {
    fmt.Printf("Running Cryptbin in developer mode\n")
  }
}
