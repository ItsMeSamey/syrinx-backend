package Server

import (
  "log"
  "strconv"
  "net/http"
  "encoding/hex"
  
  "ccs.ctf/DB"
  "ccs.ctf/GdHandler"
  
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func leaderboardHandler(c *gin.Context) {
  width, err := strconv.Atoi(c.Param("width"))
  if err != nil {
    setErrorJson(c, http.StatusBadRequest, "width parsing error\n" + err.Error())
  }
  page, err := strconv.Atoi(c.Param("page"))
  if err != nil {
    setErrorJson(c, http.StatusBadRequest, "page parsing error\n" + err.Error())
  }

  batchSize := int32(width)
  limit := int64(width)
  skip := limit*int64(page)
  cursor, err :=  DB.TeamDB.Coll.Find(DB.TeamDB.Context, bson.M{}, &options.FindOptions{
    BatchSize: &batchSize,
    Sort: bson.M{"points": 1},
    Limit: &limit,
    Skip: &skip,
  })
  if err != nil {
    setErrorJson(c, http.StatusInternalServerError, err.Error())
  }

  var teams []struct{
    N string `bson:"teamName"`
    P int    `bson:"points"`
    L int    `bson:"level"`
  }
  if err = cursor.All(DB.TeamDB.Context, &teams); err != nil {
    setErrorJson(c, http.StatusInternalServerError, err.Error())
  }
  
  c.JSON(http.StatusOK, teams)
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

  lobbyObj, err := GdHandler.LobbyIDFromUserSessionID(user.SessionID)
  if err != nil {
    setErrorJson(c, http.StatusInternalServerError, err.Error())
    return
  }

  setSuccessJson(c, gin.H{"LobbyID": hex.EncodeToString((*lobbyObj)[:])})
}

func lobbyHandler(c *gin.Context) {
  ID, err := hex.DecodeString(c.Param("lobbyID"))
  if err != nil {
    setErrorJson(c, http.StatusInternalServerError, err.Error())
    return
  }

  if len(ID) != 12 {
    setErrorJson(c, http.StatusInternalServerError, "lobbyHandler: LobbyID length mismatch")
    return
  }

  if err := GdHandler.ConnectToLobby(DB.ObjID(ID), c); err != nil {
    log.Println(err)
  }
}

