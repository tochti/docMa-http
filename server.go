package main

import (
  "path"
  "github.com/gin-gonic/gin"
  "github.com/rrawrriw/bebber"
)

func main() {
  router := gin.Default()
  htmlDir := path.Join(bebber.GetSettings("BEBBER_PUBLIC"), "html")
  router.Use(bebber.Serve("/", bebber.LocalFile(htmlDir, false)))
  router.POST("/LoadDir", bebber.LoadDir)
  router.POST("/AddTags", bebber.AddTags)
  router.Static("/public", bebber.GetSettings("BEBBER_PUBLIC"))
  serverStr := bebber.GetSettings("BEBBER_IP") +":"+
               bebber.GetSettings("BEBBER_PORT")
  router.Run(serverStr)
}

