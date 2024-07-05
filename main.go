package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
  roomID := c.Param("lobbyID")
  fmt.Fprintf(c.Writer, "Hello from room: %s", roomID)
}

func main() {
  var intt []int = nil
  fmt.Println(len(intt))

  intt = append(intt, 1)
  fmt.Println(len(intt))

  intt = append(intt[:0], intt[0])
  fmt.Println(len(intt))


  router := gin.Default()
  router.GET("/ws/:lobbyID", handler)
  router.Run(":8080")
}


