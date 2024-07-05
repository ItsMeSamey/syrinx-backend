package GdHandler

import (
  "ccs.ctf/DB"
  "github.com/gin-gonic/gin"
)


func InitLobbies(gc *gin.Context) error {
  lobbyID := gc.Param("lobbyID")
  DB.ForEachLobby(func (k, v []byte) error {
    return nil
  })
  return nil
}
