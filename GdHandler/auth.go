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

	if len(message) != 64 {
		return 0, errors.New("getUserAuth: Invalid token length!")
	}

	// Validate token
	for i, player := range lobby.Lobby.Players {
		if player.SessionID == DB.SessID(message) {
			return byte(i), nil
		}
	}

	return 0, errors.New("getUserAuth: Denied")
}

