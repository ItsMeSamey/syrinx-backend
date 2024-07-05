package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
  roomID := c.Param("roomID")
  fmt.Fprintf(c.Writer, "Hello from room: %s", roomID)
}

func main() {
  router := gin.Default()
  router.GET("/ws/:roomID", handler)
  router.Run(":8080")
}


