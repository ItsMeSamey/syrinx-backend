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

/// Handles binary message to websocket
func (lobby *Lobby) handleBinaryMessage(myIndex byte, message []byte) error {
  if len(message) < 1 {
    return errors.New("handleBinaryMessage: Unknown procedure")
  }
  procudure := message[0]
  switch (procudure) {
  case 1: //! Add player on [2, playerIndex], and send offers
    return lobby.announceToAll([]byte{1, myIndex})
  case 2: //! Remove player on [3, playerIndex]
    return lobby.announceToAll([]byte{2, myIndex})
  case 3: /// Send message to a specific person
    if len(message) < 2 {
      return errors.New("handleBinaryMessage: Cannot broadcast to Unknown")
    } else if len(message) < 3 {
      return errors.New("handleBinaryMessage: Cannot broadcast empty message")
    }

    to := message[1]
    message[1] = myIndex
    if to == myIndex {
      return lobby.announceToAll(message)
    } else {
      return lobby.announceToOne(to, message)
    }
  }
  return nil
}

/// DONOT USE THIS UNLESS YOU KNOW WHAT YOU ARE DOING. USE `announceToOne` INSTEAD
/// Send a message to someone without locking
func (lobby * Lobby) announceToOneUnlocked(Index byte, message []byte) error {
  if int(Index) > len(lobby.players) {
    return errors.New("announceToOneNOLOCK: invalid Index")
  }
  pc := lobby.players[Index].IN
  if pc != nil {
    pc <- message
  }
  return nil
}

/// Send a message to someone
func (lobby * Lobby) announceToOne(Index byte, message []byte) error {
  lobby.playerMutex.RLock()
  defer lobby.playerMutex.RUnlock()
  return lobby.announceToOneUnlocked(Index, message)
}

/// Send a message to everyone in your lobby
func (lobby * Lobby) announceToAll(message []byte) error {
  var err error = nil
  lobby.playerMutex.RLock()
  defer lobby.playerMutex.RUnlock()

  for i := range lobby.players {
    if e := lobby.announceToOneUnlocked(byte(i), message); e != nil {
      err = e
    }
  }
  return err
}

