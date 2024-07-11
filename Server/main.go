package Server

import (
  "fmt"
  "io"
  "log"
  "math"
  "os"
  "strconv"
  
  "github.com/gin-gonic/gin"
)

/// Where the log is written to
var writer io.Writer = nil
var file *os.File = nil

func handler(c *gin.Context) {
  go func(){ file.Sync() }()

  width, err := strconv.Atoi(c.Param("width"))
  if err != nil {
    fmt.Fprintf(c.Writer, "Error: " + err.Error())
    return
  } else if width > 1024*1024*32 { // 32 mb size limit
    fmt.Fprintf(c.Writer, "Error: " + strconv.Itoa(width) + " is larger than 32mb\n")
    return
  }

  pageStr := c.Param("page")
  page, err := strconv.Atoi(pageStr)
  if err != nil {
    fmt.Fprintf(c.Writer, "Error: " + err.Error())
    return
  }

  start := page*width
  if _, err = file.Seek(int64(start), io.SeekStart); err != nil {
    fmt.Fprintf(c.Writer, "Error: " + err.Error())
    return
  }

  buffer := make([]byte, width)
  _, err = file.Read(buffer)
  if _, err = file.Seek(int64(start), io.SeekStart); err != nil {
    fmt.Fprintf(c.Writer, "Error: " + err.Error())
    return
  }

  stat, err := file.Stat()
  if err != nil {
    fmt.Fprint(c.Writer, "Error getting total number of pages:  ", err.Error(), "\n", string(buffer))
  } else {
    fmt.Fprint(c.Writer, "Page ", pageStr, " of ", int64(math.Ceil( float64(stat.Size())/float64(width) )), "\n", string(buffer))
  }
}

/// This gunction Starts the frontend Server
/// This blocks forever and thus you might consider running this as async
///
/// `prepend` is the parh to `npm run build`'s output dir, usually 'dist'
func Start(ip string, prepend string) {

  file, err := os.Create("gin.log")
  if err != nil {
    log.Fatal("Could not create a log file.")
  }
  writer = io.MultiWriter(file)

  /// Disable Color to make file readable
  gin.DisableConsoleColor()

  /// Log to a file.
  gin.DefaultWriter = writer

  router := gin.Default()

  router.GET("/TYgaHxwqqKaGwMClg2fGK8CQ3rO/TOt2myVtXDYTSaBAU/R161AI3hyLEJ5R6zd3/:width/:page", handler)

  /// The signup route
  router.POST("/signup", signupHandler)

  /// The authantication route
  router.POST("/authanticate", authanticationHandler)

  router.POST("/getlobby", lobbyHandler)

  log.Println("Server satarted successfully")
  router.Run(ip)
}

