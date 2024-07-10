package Server

import (
  "github.com/gin-gonic/gin"
)

/// This gunction Starts the frontend Server
/// This blocks forever and thus you might consider running this as async
///
/// `prepend` is the parh to `npm run build`'s output dir, usually 'dist'
func Start(ip string, prepend string) {
  router := gin.Default()

  /// Files needed to serve the site
  router.StaticFile("/", prepend+"index.html")
  router.StaticFile("/bg.jpg", prepend+"bg.jpg")
  router.StaticFile("/ccs.png", prepend+"ccs.png")
  router.StaticFile("/logos.png", prepend+"logos.png")
  router.Static("/assets", prepend+"assets")

  /// The signup route
  router.POST("/signup", signupHandler)

  /// The authantication route
  router.POST("/authanticate", authanticationHandler)

  router.POST("/getlobby", lobbyHandler)

  router.Run(ip)
}

