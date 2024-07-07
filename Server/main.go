package Server

import (
  "fmt"
  "github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
  roomID := c.Param("lobbyID")
  fmt.Fprintf(c.Writer, "Hello from room: %s", roomID)
}

func Start() {
  router := gin.Default()
  // router.GET("/ws/:lobbyID", handler)
  router.POST("/signup", signupHandler)
  router.Static("/assets", "./dist/assets")
  router.StaticFile("/", "./dist/index.html")
  router.StaticFile("/bg.jpg", "./dist/bg.jpg")
  router.StaticFile("/ccs.png", "./dist/ccs.png")
  router.StaticFile("/logos.png", "./dist/logos.png")
  router.Run("127.0.0.1:8080")
}
