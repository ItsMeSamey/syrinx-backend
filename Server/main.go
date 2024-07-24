package Server

import (
  "io"
  "os"
  "log"
  
  "github.com/gin-gonic/gin"
)

var (
  /// Where the log is written to
  writer io.Writer = nil
  file *os.File = nil
)

/// This gunction Starts the frontend Server
/// This blocks forever and thus you might consider running this as async
///
/// `prepend` is the parh to `npm run build`'s output dir, usually 'dist'
func Start(ip string, prepend string) {
  initLogger()

  secret := os.Getenv("SECRET_PATH")
  if secret == "" {
    log.Fatal("Secret path not provided")
  } else if len(secret) < 8 {
    log.Fatal("Secret path is too short")
  }

  /// Logging options
  gin.DisableConsoleColor()
  gin.DefaultWriter = writer

  router := gin.Default()

  /// Logs are displayed at this route
  router.GET("/" + secret + "/:width/:page", logsHandler)

  /// Get Leaderbord Route
  router.GET("/leaderboard", leaderboardHandler)

  /// The signup route
  router.POST("/signup", signupHandler)

  /// The authantication route
  router.POST("/authanticate", authanticationHandler)

  /// Get team members
  router.POST("/teaminfo", teamInfoHandler)

  /// Get the lobbyID
  router.POST("/getlobby", getLobbyHandler)

  /// Join the lobby using WS
  router.GET("/lobby/:lobbyID", lobbyHandler)

  log.Println("Server satarted successfully")
  router.Run(ip)
}

