package main

import (
  "fmt"
  "os"
  "path"
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/rrawrriw/bebber"
  "gopkg.in/mgo.v2"
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

  config, db := SetupDefault()
  globals := bebber.Globals{Config: config, MongoDB: db}

  auth := bebber.MakeGlobalsHandler(bebber.Auth, globals)

  router := gin.Default()

  htmlDir := path.Join(config["PUBLIC_DIR"], "html")
  router.Use(bebber.Serve("/", bebber.LocalFile(htmlDir, false)))
  router.GET("/User/:name", auth, bebber.GetUser)
  router.POST("/Login", bebber.Login)
  router.Static("/public", config["PUBLIC_DIR"])
  serverStr := config["HTTP_HOST"] +":"+ config["HTTP_PORT"]

  router.Run(serverStr)

  //router.GET("/LoadBox/:boxname", bebber.Auth(), bebber.LoadBox)
  //router.POST("/AddTags", bebber.Auth(), bebber.AddTags)
  //router.GET("/LoadFile/:boxname/:filename", bebber.Auth(), bebber.LoadFile)
  //router.POST("/MoveFile", bebber.Auth(), bebber.MoveFile)
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


func SetupDefault() (bebber.Config, bebber.MongoDBConn) {
  config := bebber.Config{}
  config["FILES_DIR"] = bebber.GetSettings("BEBBER_FILES")
  config["PUBLIC_DIR"] = bebber.GetSettings("BEBBER_PUBLIC")
  config["HTTP_HOST"] = bebber.GetSettings("BEBBER_IP")
  config["HTTP_IP"] = bebber.GetSettings("BEBBER_PORT")
  config["MONGODB_HOST"] = bebber.GetSettings("BEBBER_DB_SERVER")
  config["MONGODB_DBNAME"] = bebber.GetSettings("BEBBER_DB_NAME")

  dialInfo := &mgo.DialInfo{
                Addrs: []string{config["MONGODB_HOST"]},
              }
  session, err := mgo.DialWithInfo(dialInfo)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(2)
  }

  conn := bebber.MongoDBConn {
            DialInfo: dialInfo,
            Session: session,
            DBName: config["MONGODB_DBNAME"],
          }

  return config, conn
}
