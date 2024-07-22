package Server

import (
  "io"
  "os"
  "fmt"
  "log"
  "math"
  "errors"
  "strconv"
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

func handleLogs(c *gin.Context) {
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

