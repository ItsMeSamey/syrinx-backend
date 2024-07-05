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

func (lobby *Lobby) getUserAuth(messageType int, message []byte) (*DB.User, error) {
	if (messageType != websocket.TextMessage) {
		return nil, errors.New("getUserAuth: Invalid messageType")
	}
	data := message[1:]
	var user *DB.User
	var err error
	if message[0] == '0' {
		user, err = DB.UserFromSessionID(data)
	} else if message[0] == '1' {
		var auth authUser
		err := json.Unmarshal(message, &auth)
		if err != nil {
			return nil, errors.New("getUserAuth: Json Unmarshal Error")
		}
		user, err = DB.UserAuthenticate(auth.Username, auth.Password)
	} else {
		return nil, errors.New("getUserAuth: spec violation")
	}

	has, err := user.UserInLobby(lobby.ID)
	if err != nil {
		return nil, err
	}
	if has != true {
		return nil, errors.New("getUserAuth: Policy Violation")
	}
	return user, err
}

