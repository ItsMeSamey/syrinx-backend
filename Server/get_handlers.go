package Server

import (
  "io"
  "log"
  "fmt"
  "math"
  "strconv"
  "net/http"
  "encoding/hex"
  
  "ccs.ctf/DB"
  "ccs.ctf/GdHandler"
  
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func lobbyHandler(c *gin.Context) {
  ID, err := hex.DecodeString(c.Param("lobbyID"))
  if err != nil {
    setErrorJson(c, http.StatusInternalServerError, err.Error())
    return
  }

  if len(ID) != 3 {
    setErrorJson(c, http.StatusInternalServerError, "lobbyHandler: LobbyID length mismatch")
    return
  }

  if err := GdHandler.ConnectToLobby(DB.TID(ID), c); err != nil {
    log.Println(err)
  }
}

func leaderboardHandler(c *gin.Context) {
  batchSize := int32(1024)
  limit := int64(1024)
  cursor, err :=  DB.TeamDB.Coll.Find(DB.TeamDB.Context, bson.M{}, &options.FindOptions{
    BatchSize: &batchSize,
    Sort: bson.M{"points": -1},
    Limit: &limit,
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

func logsHandler(c *gin.Context) {
  go func(){ file.Sync() }()

  width, err := strconv.Atoi(c.Param("width"))
  if err != nil {
    fmt.Fprintf(c.Writer, "Error: Width parsing error\n" + err.Error())
    return
  } else if width > 1024*1024*32 { // 32 mb size limit
    fmt.Fprintf(c.Writer, "Error: Width " + strconv.Itoa(width) + " is larger than 32mb")
    return
  }

  pageStr := c.Param("page")
  page, err := strconv.Atoi(pageStr)
  if err != nil {
    fmt.Fprintf(c.Writer, "Error: Page parsing error\n" + err.Error())
    return
  }

  if _, err = file.Seek(int64(page*width), io.SeekStart); err != nil {
    fmt.Fprintf(c.Writer, "Error: Seek error\n" + err.Error())
    return
  }

  buffer := make([]byte, width)
  n, err := file.Read(buffer)
  if err != nil {
    fmt.Fprintf(c.Writer, "Error: Read error\n" + err.Error())
    return
  }
  buffer = buffer[:n]

  stat, err := file.Stat()
  if err != nil {
    fmt.Fprint(c.Writer, "Error: getting file stats failed", err.Error(), "\n", string(buffer))
  } else {
    fmt.Fprint(c.Writer, "Page ", pageStr, "/", int64(math.Ceil( float64(stat.Size())/float64(width) )) - 1, "\n", string(buffer))
  }
}

