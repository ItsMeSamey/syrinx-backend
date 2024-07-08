package main

import (
  "log"

  "ccs.ctf/DB"
  "ccs.ctf/Server"
)

func main() {
  // Connect to the db
  err := DB.InitDB("mongodb://localhost:27017")
  if err != nil {
    log.Fatal(err)
  }

  // Start serving the front end server
  Server.Start("127.0.0.1:8080", "../Syrinx_Login/dist/")
}

