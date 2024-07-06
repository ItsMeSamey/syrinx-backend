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
	data := string(message[1:])
	var user *DB.User = nil
	var err error = nil
	if message[0] == '0' {
		// We got a _id
		user, err = DB.UserFromSessionID(data)
	} else if message[0] == '1' {
		// We got Username and pass
		var auth authUser
		err := json.Unmarshal(message, &auth)
		if err != nil {
			return nil, errors.New("getUserAuth: Json Unmarshal Error")
		}
		user, err = DB.UserAuthenticate(auth.Username, auth.Password)
	} else {
		// Unknown auth type
		return nil, errors.New("getUserAuth: spec violation")
	}
	if err != nil {
		return nil, err
	}

	// check if the user is actually in lobby by lobbyID
	has, err := user.UserInLobby(lobby.ID)
	if err != nil {
		return nil, err
	}
	if has != true {
		return nil, errors.New("getUserAuth: Policy Violation")
	}
	return user, nil
}

