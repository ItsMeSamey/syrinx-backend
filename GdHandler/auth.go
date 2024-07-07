package GdHandler

import (
	"errors"

	"ccs.ctf/DB"

	"github.com/gorilla/websocket"
)

func (lobby *Lobby) getUserAuth(messageType int, message []byte) (string, error) {
	if (messageType != websocket.TextMessage) {
		return "", errors.New("getUserAuth: Invalid messageType")
	}

	data := message[1:]
	var user *DB.User = nil
	var err error = nil
	if message[0] == '0' {
		// We got a SessID
		user, err = DB.UserFromSessionID(DB.SessID(data))
	} else {
		// Unknown auth type
		return "", errors.New("getUserAuth: Unknown Proto")
	}
	if err != nil {
		return "", err
	}

	// check if the user is has permission to be in lobby
	has, err := DB.UserInLobby(lobby.ID, user.ID)
	if err != nil {
		return "", err
	}
	if has != true {
		return "", errors.New("getUserAuth: Denied")
	}
	return user.ID, nil
}

