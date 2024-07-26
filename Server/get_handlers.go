package Server

import (
  "io"
  "fmt"
  "math"
  "strconv"
  "net/http"
  "encoding/hex"
  
  "ccs.ctf/DB"
  "ccs.ctf/GdHandler"
  
  "github.com/gin-gonic/gin"
)

func lobbyHandler(c *gin.Context) {
  if !DB.State.GameOn {
    setErrorJson(c, http.StatusBadRequest, "GAME STOPPED!!")
    return
  }

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
    // fmt.Println(err)
  }
}

func leaderboardHandler(c *gin.Context) {
  type TINFO struct{
    I DB.TID
    N string
    P int
    L int
  }
  var teams []TINFO

  GdHandler.TEAMSMutex.RLock()
  for _, val := range GdHandler.TEAMS {
    teams = append(teams, TINFO{
      I: val.TeamID,
      N: val.TeamName,
      P: val.Points,
      L: val.Level,
    })
  }
  GdHandler.TEAMSMutex.RUnlock()

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

