package main

import (
  "log"

  "ccs.ctf/DB"
  "ccs.ctf/Server"
  "ccs.ctf/GdHandler"
)

func main() {
  // Connect to the db
  if err := DB.Init(); err != nil {
    log.Fatal(err)
  }

  // Connect to the db
  if err := GdHandler.Init(); err != nil {
    log.Fatal(err)
  }

  // Start serving the front end server
  Server.Start("0.0.0.0:8080", "../Syrinx_Login/dist/")
}

