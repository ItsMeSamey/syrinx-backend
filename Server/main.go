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

/// This gunction Starts the frontend Server
/// This blocks forever and thus you might consider running this as async
///
/// `prepend` is the parh to `npm run build`'s output dir, usually 'dist'
func Start(ip string, prepend string) {
  var err error

  secret := os.Getenv("SECRET_PATH")
  if secret == "" {
     log.Fatal("Secret path not provided")
  } else if len(secret) < 8 {
    log.Fatal("Secret path is too short")
  }

  file, err = os.OpenFile("./log/gin.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
  if err != nil {
    log.Fatal("Could not create a log file.")
  }
  if os.Getenv("GIN_MODE") != "release" {
    writer = io.MultiWriter(file, os.Stdout)
  } else {
    writer = io.Writer(file)
  }

  /// Disable Color to make file readable
  gin.DisableConsoleColor()

  /// Log to a file.
  gin.DefaultWriter = writer

  router := gin.Default()

  router.GET("/" + secret + "/:width/:page", handler)

  /// The signup route
  router.POST("/signup", signupHandler)

  /// The authantication route
  router.POST("/authanticate", authanticationHandler)

  router.POST("/getlobby", lobbyHandler)

  log.Println("Server satarted successfully")
  router.Run(ip)
}

