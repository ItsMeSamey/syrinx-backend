package Server

import (
  "net/http"
  "encoding/hex"
  
  "ccs.ctf/DB"
  "ccs.ctf/GdHandler"
  
  "github.com/gin-gonic/gin"
)

/// Function to call on signup request
func signupHandler(c *gin.Context) {
  var user DB.CreatableUser
  if err := bindJson(c, &user); err != nil {
    setErrorJson(c, http.StatusBadRequest, err.Error())
    return
  }

  SessionID, err := DB.CreateUser(&user)
  if err != nil {
    setErrorJson(c, http.StatusInternalServerError, err.Error())
  } else {
    setSuccessJson(c, gin.H{"SessionID": SessionID, "TeamID": user.TeamID})
  }
}

/// Function to call for authantication
func authanticationHandler(c *gin.Context) {
  var user DB.User
  if err := bindJson(c, &user); err != nil {
    setErrorJson(c, http.StatusBadRequest, err.Error())
    return
  }

  if user.Username == "" || user.Password == "" {
    setErrorJson(c, http.StatusBadRequest, "Username and Password are required")
    return
  }

  usr, err := DB.UserAuthenticate(user.Username, user.Password)
  if err != nil {
    setErrorJson(c, http.StatusUnauthorized, err.Error())
    return
  }

  setSuccessJson(c, gin.H{"SessionID": usr.SessionID})
}

/// Function to call when user asks for their lobby
func getLobbyHandler(c *gin.Context) {
  var user struct { SessionID DB.SessID }

  if err := bindJson(c, &user); err != nil {
    setErrorJson(c, http.StatusBadRequest, err.Error())
    return
  }

  if user.SessionID == nil {
    setErrorJson(c, http.StatusBadRequest, "SessionID is required")
    return
  }

  lobbyObj, level, err := GdHandler.LobbyIDFromUserSessionID(user.SessionID)
  if err != nil {
    setErrorJson(c, http.StatusInternalServerError, err.Error())
    return
  }

  setSuccessJson(c, gin.H{"LobbyID": hex.EncodeToString((*lobbyObj)[:]), "Level": level, })
}

func teamInfoHandler(c *gin.Context) {
  var userID struct {SessionID DB.SessID}

  if err := bindJson(c, &userID); err != nil {
    setErrorJson(c, http.StatusBadRequest, err.Error())
    return
  }

  if userID.SessionID == nil {
    setErrorJson(c, http.StatusBadRequest, "SessionID is required")
    return
  }

  user, err := DB.UserFromSessionID(userID.SessionID)
  if err != nil {
    setErrorJson(c, http.StatusBadRequest, err.Error())
    return
  }

  all, team, err := GdHandler.GetTeamAndPlayers(user.TeamID)
  if err != nil {
    setErrorJson(c, http.StatusInternalServerError, err.Error())
    return
  }

  setSuccessJson(c, gin.H{
    "T": user.TeamID,
    "M": user.Username,
    "A": all,
    "N": team.TeamName,
    "P": team.Points,
  })
}

