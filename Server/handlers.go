package Server

import (
  "encoding/json"
  "errors"
  "io"
  "fmt"
  "net/http"
  "time"
  "encoding/base64"

  "ccs.ctf/DB"
  "ccs.ctf/GdHandler"
  
  "github.com/gin-gonic/gin"
)

func bindJson(c *gin.Context, obj any) error {
  jsonData, err := io.ReadAll(c.Request.Body)
  if err != nil {
    return errors.New("Server.bindJson error: \n" + err.Error())
  }

  /// Logging code
  writer.Write([]byte("\n\n>>>>>>>>>>" + time.Now().String() + "\n>> body\n"))
  writer.Write(jsonData)
  writer.Write([]byte("\n << body\n"))
  /// Logging

  err = json.Unmarshal(jsonData, obj);
  if err != nil {
    return errors.New("Server.bindJson error: \n" + err.Error())
  }
  return nil
}

/// Easily extensible for logging
func setJson(c *gin.Context, code int, json gin.H) {
  /// Logging code
  writer.Write([]byte(">> response"))
  for key, value := range json {
    writer.Write([]byte(key))
    writer.Write([]byte(": "))
    switch t := value.(type) {
    case []byte:
      writer.Write([]byte(base64.StdEncoding.EncodeToString(t)))
    default:
      writer.Write([]byte(fmt.Sprintf("%v", value)))
    }
    writer.Write([]byte("\n"))
  }
  writer.Write([]byte("<< response\n<<<<<"))
  /// Logging

  c.JSON(code, json)
}

func setSuccessJson(c *gin.Context, json gin.H) {
  setJson(c, http.StatusOK, json)
}

func setErrorJson(c *gin.Context, code int, errstr string) {
  // go func(){ _, _ = writer.Write([]byte(errstr)) }()
  setJson(c, code, gin.H{"error": errstr})
}

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
func lobbyHandler(c *gin.Context) {
  var user struct { SessionID DB.SessID }

  if err := bindJson(c, &user); err != nil {
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

