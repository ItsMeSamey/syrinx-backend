package Server

import (
  "fmt"
  "github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
  roomID := c.Param("lobbyID")
  fmt.Fprintf(c.Writer, "Hello from room: %s", roomID)
}

func Start(prepend string) {
  router := gin.Default()

  /// Files needed to serve the site
  router.StaticFile("/", prepend+"index.html")
  router.StaticFile("/bg.jpg", prepend+"bg.jpg")
  router.StaticFile("/ccs.png", prepend+"ccs.png")
  router.StaticFile("/logos.png", prepend+"logos.png")
  router.Static("/assets", prepend+"assets")

  /// The signup handler
  router.POST("/signup", signupHandler)

  // router.GET("/ws/:lobbyID", handler)

  router.Run("127.0.0.1:8080")
}
