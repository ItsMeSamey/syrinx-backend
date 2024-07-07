package GdHandler

import (
	"errors"

	"ccs.ctf/DB"

	"github.com/gorilla/websocket"
)


/// Checks user auth
func (lobby *Lobby) getUserAuth(messageType int, message []byte) (byte, error) {
	if (messageType != websocket.BinaryMessage) {
		return 0, errors.New("getUserAuth: Invalid messageType")
	}

	// Validate token
	user, err := DB.UserFromSessionID(DB.SessID(message))
	if err != nil {
		return 0, err
	}

	// Checks if the user is a member of this lobby
	for i, player := range lobby.players {
		if player.ID == user.ID {
			return byte(i), nil
		}
	}

	return 0, errors.New("getUserAuth: Denied")
}

