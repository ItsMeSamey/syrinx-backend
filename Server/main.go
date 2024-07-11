package Server

import (
  "io"
  "log"
  "os"
  
  "github.com/gin-gonic/gin"
)

/// Where the log is written to
var writer io.Writer = nil

/// This gunction Starts the frontend Server
/// This blocks forever and thus you might consider running this as async
///
/// `prepend` is the parh to `npm run build`'s output dir, usually 'dist'
func Start(ip string, prepend string) {

  f, err := os.Create("gin.log")
  if err != nil {
    log.Fatal("Could not create a log file.")
  }
  writer = io.MultiWriter(f)

  /// Disable Color to make file readable
  gin.DisableConsoleColor()

  /// Log to a file.
  gin.DefaultWriter = writer

  router := gin.Default()

  /// The signup route
  router.POST("/signup", signupHandler)

  /// The authantication route
  router.POST("/authanticate", authanticationHandler)

  router.POST("/getlobby", lobbyHandler)

  log.Println("Server satarted successfully")
  router.Run(ip)
}

