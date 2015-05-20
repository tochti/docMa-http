package main

import (
  "path"
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/rrawrriw/bebber"
)

func main() {
  acc := flag.Bool("steuerberater", false, "Start Server für Steuerberater")
  flag.Parse()
  if *acc {
    StartAccServer()
  } else {
    StartDefaultServer()
  }
}

func StartDefaultServer () {
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

func StartAccServer() {
  router := gin.Default()
  htmlDir := path.Join(bebber.GetSettings("BEBBER_PUBLIC"), "html")
  router.Use(bebber.Serve("/", bebber.LocalFile(htmlDir, false)))
  router.GET("/LoadAccFiles", bebber.LoadAccFiles)
  router.Static("/public", bebber.GetSettings("BEBBER_PUBLIC"))
  router.Static("/data", bebber.GetSettings("BEBBER_ACC_DATA"))
  serverStr := bebber.GetSettings("BEBBER_IP") +":"+
               bebber.GetSettings("BEBBER_PORT")
  router.Run(serverStr)
}

