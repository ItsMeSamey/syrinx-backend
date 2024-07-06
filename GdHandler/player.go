package GdHandler

import (
  "github.com/gorilla/websocket"
)

type Player struct {
  ID string
  WS *websocket.Conn
}

