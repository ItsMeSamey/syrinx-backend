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

func getUserAuth(messageType int, message []byte) (*DB.User, error) {
	if (messageType != websocket.TextMessage) {
		return nil, errors.New("getUserAuth: Invalid messageType")
	}
	data := message[1:]
	if message[0] == '0' {
		return DB.UserFromSessionID(data)
	} else if message[0] == '1' {
		var auth authUser
		err := json.Unmarshal(message, &auth)
		if err != nil {
			return nil, errors.New("getUserAuth: Json Unmarshal Error")
		}
		return DB.UserAuthenticate(auth.Username, auth.Password)
	}
	return nil, errors.New("getUserAuth: spec violation")
}

