package main

import (
  "github.com/gin-gonic/gin"
  "github.com/rrawrriw/bebber"
)

func main() {
  router := gin.Default()
  router.POST("/LoadDir", bebber.LoadDir)
  router.POST("/AddTags", bebber.AddTags)
  serverStr := bebber.GetSettings("BEBBER_IP") +":"+
               bebber.GetSettings("BEBBER_PORT")
  router.Run(serverStr)
}

