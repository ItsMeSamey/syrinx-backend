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
  err = DB.SendConfirmationEmail(&DB.User{
    Username: "a",
    Email: "sanyamsingh504@gmail.com",
    Password: "adsasdasd",
    DiscordID: "asdasddas",
    TeamID: DB.TID(&[3]byte{1,2,3}),
  })
  if err != nil {
    log.Fatal(err)
  }

  Server.Start("127.0.0.1:8080", "../Syrinx_Login/dist/")
}

