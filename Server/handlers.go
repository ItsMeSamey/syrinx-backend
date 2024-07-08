package Server

import (
  "net/http"
  
  "ccs.ctf/DB"
  "ccs.ctf/GdHandler"
  
  "github.com/gin-gonic/gin"
)

func signupHandler(c *gin.Context) {
  var user DB.User
  if err := c.BindJSON(&user); err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  err :=  DB.CreateUser(&user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  c.JSON(http.StatusOK, gin.H{"SessionID": user.SessionID, "TeamID": user.TeamID})
}

func authanticationHandler(c *gin.Context) {
  var user DB.User
  if err := c.BindJSON(&user); err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  if user.Username == "" || user.Password == "" {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username and Password are required"})
    return
  }

  usr, err := DB.UserAuthenticate(user.Username, user.Password)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    return
  }

  c.JSON(http.StatusOK, gin.H{"SessionID": usr.SessionID})
}

func lobbyHandler(c *gin.Context) {
  var user DB.User
  if err := c.BindJSON(&user); err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  if user.SessionID == nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username and Password are required"})
    return
  }

  lobbyObj, err := DB.LobbyFromUserSessionID(user.SessionID)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  GdHandler.ConnectToLobby(lobbyObj, c)
}

