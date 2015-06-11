package main

import (
  "path"
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/rrawrriw/bebber"
)

func main() {
  acc := flag.Bool("steuerberater", false, "Start Server für Steuerberater")
  valid := flag.Bool("valid", false, "Prüft ob für alle Buchungsätze eine Datei vorhanden ist")
  flag.Parse()
  if *acc && *valid{
    StartAccServer(true)
  } else if *acc {
    StartAccServer(false)
  } else {
    StartDefaultServer()
  }
}

func StartDefaultServer () {
  router := gin.Default()
  htmlDir := path.Join(bebber.GetSettings("BEBBER_PUBLIC"), "html")
  router.Use(bebber.Serve("/", bebber.LocalFile(htmlDir, false)))
  router.GET("/LoadBox/:boxname", bebber.Auth(), bebber.LoadBox)
  router.POST("/AddTags", bebber.Auth(), bebber.AddTags)
  router.GET("/User/:name", bebber.Auth(), bebber.GetUser)
  router.POST("/MoveFile", bebber.Auth(), bebber.MoveFile)
  router.POST("/Login", bebber.Login)
  router.Static("/public", bebber.GetSettings("BEBBER_PUBLIC"))
  serverStr := bebber.GetSettings("BEBBER_IP") +":"+
               bebber.GetSettings("BEBBER_PORT")
  router.Run(serverStr)
}

func StartAccServer(valid bool) {
  router := gin.Default()
  htmlDir := path.Join(bebber.GetSettings("BEBBER_PUBLIC"), "html")
  router.Use(validCSV(valid))
  router.Use(bebber.Serve("/", bebber.LocalFile(htmlDir, false)))
  router.GET("/LoadAccFiles", bebber.LoadAccFiles)
  router.Static("/public", bebber.GetSettings("BEBBER_PUBLIC"))
  router.Static("/data", bebber.GetSettings("BEBBER_ACC_DATA"))
  serverStr := bebber.GetSettings("BEBBER_IP") +":"+
               bebber.GetSettings("BEBBER_PORT")
  router.Run(serverStr)
}

func validCSV(valid bool) gin.HandlerFunc {
  return func (c *gin.Context) {
    c.Set("validCSV", valid)
    c.Next()
  }
}

