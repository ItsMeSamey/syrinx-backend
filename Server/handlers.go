package Server

import (
  "io"
  "net/http"
  
  "ccs.ctf/DB"
  "ccs.ctf/GdHandler"
  
  "github.com/gin-gonic/gin"
)

/// Easily extensible for logging pur
func setJson(c *gin.Context, code int, json gin.H) {
  go func(){
    body, err := c.Request.GetBody()
    if err != nil || body == nil { return }
    data, err := io.ReadAll(body)
    if err != nil || data == nil { return }
    _, _ = writer.Write(data)
  }()
  c.JSON(code, json)
}

func setSuccessJson(c *gin.Context, json gin.H) {
  setJson(c, http.StatusOK, json)
}

func setErrorJson(c *gin.Context, code int, errstr string) {
  go func(){ writer.Write([]byte(errstr)) }()
  setJson(c, code, gin.H{"error": errstr})
}

/// Function to call on signup request
func signupHandler(c *gin.Context) {
  var user DB.User
  if err := c.BindJSON(&user); err != nil {
    setErrorJson(c, http.StatusBadRequest, err.Error())
    return
  }

  if  err := DB.CreateUser(&user); err != nil {
    setErrorJson(c, http.StatusBadRequest, err.Error())
    return
  }

  setSuccessJson(c, gin.H{"SessionID": user.SessionID, "TeamID": user.TeamID})
}

/// Function to call for authantication
func authanticationHandler(c *gin.Context) {
  var user DB.User
  if err := c.BindJSON(&user); err != nil {
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
func lobbyHandler(c *gin.Context) {
  var user struct {

  }
  if err := c.BindJSON(&user); err != nil {
    setErrorJson(c, http.StatusBadRequest, err.Error())
    return
  }

  if user.SessionID == nil {
    setErrorJson(c, http.StatusBadRequest, "Username and Password are required")
    return
  }

  lobbyObj, err := DB.LobbyFromUserSessionID(user.SessionID)
  if err != nil {
    setErrorJson(c, http.StatusInternalServerError, err.Error())
    return
  }

  GdHandler.ConnectToLobby(lobbyObj, c)
}

