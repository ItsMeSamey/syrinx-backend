package GdHandler

import (
	"encoding/json"
	"errors"
	"log"

	"ccs.ctf/DB"
	"github.com/gorilla/websocket"
)

/// This will probably handle questioning/answering
func (lobby *Lobby) handleTextMessage(myIndex byte, message []byte, conn *websocket.Conn) error {
  log.Println("Got String: ", message)
  _question := DB.Question{}
  if err := json.Unmarshal(message, _question); err != nil {
    return err
  }

  if _question.Answer == "" {
    question, err := DB.QuestionFromID(_question.ID)
    if err != nil {
      return err
    }
    conn.WriteMessage(websocket.TextMessage, []byte(question.Question))
  } else {
    myTeam := lobby.Lobby.Players[myIndex].TeamID
    for i := range lobby.Lobby.Teams {
      if *(lobby.Lobby.Teams[i].TeamID) == *myTeam {
        lobby.PlayerMutex.Lock()
        lobby.Lobby.Teams[i].CheckAnswer(_question.ID, _question.Answer, 5)
        lobby.Lobby.Teams[i].SyncTryHard(5)
        lobby.PlayerMutex.Unlock()
        break
      }
    }
  }

  return nil 
}

/// Handles binary message to websocket
func (lobby *Lobby) handleBinaryMessage(myIndex byte, message []byte) error {
  if len(message) < 1 {
    return errors.New("handleBinaryMessage: Error empty message")
  }
  procudure := message[0]
  switch (procudure) {
  case 1: //! Add player on [2, playerIndex], and send offers
    log.Println("Got Binary: ", message)
    return lobby.announceToAll(myIndex, []byte{0x01, myIndex})
  case 2: //! Remove player on [3, playerIndex]
    log.Println("Got Binary: ", message)
    return lobby.announceToAll(myIndex, []byte{0x02, myIndex})
  case 3: /// Send message to a specific person
    if len(message) < 2 {
      return errors.New("handleBinaryMessage: Cannot broadcast to Unknown")
    } else if len(message) < 3 {
      return errors.New("handleBinaryMessage: Cannot broadcast empty message")
    }
    log.Println("Got: ", message[:2], " ", string(message[2]))

    to := message[1]
    if to == myIndex {
      return lobby.announceToAll(myIndex, message)
    } else {
      message[1] = myIndex
      return lobby.announceToOne(to, message)
    }
  default:
    return errors.New("handleBinaryMessage: Unknown procedure")
  }
}

/// DONOT USE THIS UNLESS YOU KNOW WHAT YOU ARE DOING. USE `announceToOne` INSTEAD
/// Send a message to someone without locking
func (lobby * Lobby) announceToOneUnlocked(Index byte, message []byte) error {
  ind := int(Index)
  if ind > len(lobby.Lobby.Players) {
    return errors.New("announceToOneNOLOCK: invalid Index")
  }
  pc := lobby.Lobby.Players[ind].IN
  if pc != nil {
    pc <- message
  }
  return nil
}

/// Send a message to someone
func (lobby * Lobby) announceToOne(Index byte, message []byte) error {
  lobby.PlayerMutex.RLock()
  defer lobby.PlayerMutex.RUnlock()
  return lobby.announceToOneUnlocked(Index, message)
}

/// Send a message to everyone in your lobby
func (lobby * Lobby) announceToAll(myIndex byte, message []byte) error {
  var err error = nil
  lobby.PlayerMutex.RLock()
  defer lobby.PlayerMutex.RUnlock()

  for i := range lobby.Lobby.Players {
    if i != int(myIndex){
      err = errors.Join(err, lobby.announceToOneUnlocked(byte(i), message))
    }
  }

  if err != nil {
    return errors.New("announceToAll: announceToOne error\n" + err.Error())
  }
  return nil
}

