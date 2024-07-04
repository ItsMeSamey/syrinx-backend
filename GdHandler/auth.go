package GdHandler

import (
	"encoding/json"
	"errors"

	"ccs.ctf/DB"

	"github.com/gorilla/websocket"
	bolt "go.etcd.io/bbolt"
)

type authUser struct {
	username string `json:"user"`
	password string `json:"pass"`
}

func getUserAuth(messageType int, message []byte) (*DB.User, error) {
	if (messageType != websocket.TextMessage) {
		return nil, errors.New("getUserAuth: Invalid messageType")
	}
	var user *DB.User
	data := message[1:]
	if (message[0] == '0') {
		DB.GetUserFromSessionID(string(data))
	}
	return &user, nil
}

