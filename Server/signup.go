package Server

import (
  "net/http"

  "ccs.ctf/DB"

  "github.com/gin-gonic/gin"
)

func signupHandler(c *gin.Context) {
  var user DB.User
  // var data map[string]interface{}
  if err := c.BindJSON(&user); err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  id, err :=  DB.CreateUser(&user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  c.JSON(http.StatusOK, gin.H{"SessionID": id})
}

