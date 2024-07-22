package Server

import (
  "io"
  "os"
  "fmt"
  "log"
  "errors"
  "net/http"
  "encoding/json"
  "encoding/base64"

  "github.com/gin-gonic/gin"
)

/// Where the log is written to
var (
  writer io.Writer = nil
  file *os.File = nil
)

func initLogger() {
  var err error

  file, err = os.OpenFile("./log/gin.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
  if err != nil {
    log.Fatal("Could not create a log file.")
  }
  if os.Getenv("GIN_MODE") != "release" {
    writer = io.MultiWriter(file, os.Stdout)
  } else {
    writer = io.Writer(file)
  }
}

func bindJson(c *gin.Context, obj any) error {
  jsonData, err := io.ReadAll(c.Request.Body)
  if err != nil {
    return errors.New("Server.bindJson error: \n" + err.Error())
  }

  /// Logging code
  writer.Write([]byte("\n>>>>>\n>> "))
  writer.Write(jsonData)
  writer.Write([]byte("\n<< "))
  /// Logging

  err = json.Unmarshal(jsonData, obj);
  if err != nil {
    return errors.New("Server.bindJson error: \n" + err.Error())
  }
  return nil
}

/// Easily extensible for logging
func setJson(c *gin.Context, code int, json gin.H) {
  /// Logging code
  for key, value := range json {
    writer.Write([]byte(key))
    writer.Write([]byte(": "))
    switch t := value.(type) {
    case []byte:
      writer.Write([]byte(base64.StdEncoding.EncodeToString(t)))
    default:
      writer.Write([]byte(fmt.Sprintf("%v", value)))
    }
    writer.Write([]byte("\n"))
  }
  writer.Write([]byte("<<<<< "))
  /// Logging

  c.JSON(code, json)
}

func setSuccessJson(c *gin.Context, json gin.H) {
  setJson(c, http.StatusOK, json)
}

func setErrorJson(c *gin.Context, code int, errstr string) {
  // go func(){ _, _ = writer.Write([]byte(errstr)) }()
  setJson(c, code, gin.H{"error": errstr})
}

