package GdHandler

import (
  "errors"
  "reflect"

  "github.com/gorilla/websocket"
)

/// Checks user auth
func (lobby *Lobby) getUserAuth(messageType int, message []byte) (byte, error) {
  if (messageType != websocket.BinaryMessage) {
    return 0, errors.New("getUserAuth: Invalid messageType")
  }

  if len(message) != 64 {
    return 0, errors.New("getUserAuth: Invalid token length!")
  }

  // Validate token
  for i, player := range lobby.Lobby.Players {
    if reflect.DeepEqual(*player.SessionID, [64]byte(message)) {
      return byte(i), nil
    }
  }

  return 0, errors.New("getUserAuth: Denied")
}

