package main

import (
  "log"
  "os"

  "ccs.ctf/DB"
  "ccs.ctf/Server"
)

func main() {
  // Connect to the db
  err := DB.InitDB(os.Getenv("MONGOURI"))
  if err != nil {
    log.Fatal(err)
  }

  // Start serving the front end server
  Server.Start("0.0.0.0:8080", "../Syrinx_Login/dist/")
}

