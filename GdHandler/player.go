package GdHandler

import (
  "ccs.ctf/DB"

  "github.com/gorilla/websocket"
)

type Player struct {
  user *DB.User
  conn *websocket.Conn
}

