package GdHandler

import (
  "errors"
)

/// This will probably handle questioning/answering
func (lobby *Lobby) handleTextMessage(myIndex byte, message []byte) error {
  // TODO: implement
  _ = myIndex
  _ = message
  return nil
}

/// Handles ninary message to websocket
func (lobby *Lobby) handleBinaryMessage(myIndex byte, message []byte) error {
  if len(message) < 1 {
    return errors.New("handleBinaryMessage: Destination Index not provided")
  }
  dest := message[0]
  lobby.playerMutex.RLock()
  defer lobby.playerMutex.RUnlock()
  if int(dest) >= len(lobby.players) {
    return errors.New("handleBinaryMessage: invalid destination")
  }
  receiver := &lobby.players[dest]
  if receiver.IN == nil {
    return errors.New("handleBinaryMessage: receiver is not connected")
  }
  message[0] = myIndex
  // FIXME: prevent deadlock somehow ?
  receiver.IN <- message
  return nil
}

