package GdHandler

import (
	"encoding/json"
	"errors"

	"ccs.ctf/DB"

	"github.com/gorilla/websocket"
)

type authUser struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

func (lobby *Lobby) getUserAuth(messageType int, message []byte) (string, error) {
	if (messageType != websocket.TextMessage) {
		return "", errors.New("getUserAuth: Invalid messageType")
	}

	data := message[1:]
	var user *DB.User = nil
	var err error = nil
	if message[0] == '0' {
		// We got a _id
		user, err = DB.UserFromSessionID(DB.SessID(data))
	} else if message[0] == '1' {
		// We got Username and pass
		var auth authUser
		err := json.Unmarshal(data, &auth)
		if err != nil {
			return "", errors.New("getUserAuth: Json Unmarshal Error")
		}
		user, err = DB.UserAuthenticate(auth.Username, auth.Password)
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

