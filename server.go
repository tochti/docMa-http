package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/tochti/docMa-handler"
	"github.com/tochti/docMa-handler/accountingData"
	"github.com/tochti/docMa-handler/common"
	"github.com/tochti/docMa-handler/dbVars"
	"github.com/tochti/docMa-handler/docNumberProposal"
	"github.com/tochti/docMa-handler/docs"
	"github.com/tochti/docMa-handler/labels"
	"github.com/tochti/gin-gum/gumauth"
	"github.com/tochti/gin-gum/gumspecs"
	"github.com/tochti/gin-gum/gumwrap"
	"github.com/tochti/smem"
	"gopkg.in/mgo.v2"
)

func main() {
	acc := flag.Bool("steuerberater", false, "Start Server für Steuerberater")
	valid := flag.Bool("valid", false, "Prüft ob für alle Buchungsätze eine Datei vorhanden ist")
	flag.Parse()
	if *acc && *valid {
		StartAccServer(true)
	} else if *acc {
		StartAccServer(false)
	} else {
		StartDefaultServer()
	}
}

func StartDefaultServer() {

	/*
		config, db := SetupDefault()
		globals := bebber.Globals{Config: config, MongoDB: db}

		makeGlobalsHandler := bebber.MakeGlobalsHandler
		authHandler := makeGlobalsHandler(bebber.Auth, globals)
		loginHandler := makeGlobalsHandler(bebber.LoginHandler, globals)
		searchDocsHandler := makeGlobalsHandler(bebber.SearchDocsHandler, globals)
		userHandler := makeGlobalsHandler(bebber.UserHandler, globals)
		docMakeHandler := makeGlobalsHandler(bebber.DocMakeHandler, globals)
		docReadHandler := makeGlobalsHandler(bebber.DocReadHandler, globals)
		docChangeHandler := makeGlobalsHandler(bebber.DocChangeHandler, globals)
		docRemoveHandler := makeGlobalsHandler(bebber.DocRemoveHandler, globals)
		docRenameHandler := makeGlobalsHandler(bebber.DocRenameHandler, globals)
		docAppendLabelsHandler := makeGlobalsHandler(bebber.DocAppendLabelsHandler, globals)
		docRemoveLabelsHandler := makeGlobalsHandler(bebber.DocRemoveLabelHandler, globals)
		docAppendDocNumbersHandler := makeGlobalsHandler(bebber.DocAppendDocNumbersHandler, globals)
		docRemoveDocNumberHandler := makeGlobalsHandler(bebber.DocRemoveDocNumberHandler, globals)
		accProcessMakeHandler := makeGlobalsHandler(bebber.AccProcessMakeHandler, globals)
		accProcessFindByDocNumberHandler := makeGlobalsHandler(bebber.AccProcessFindByDocNumberHandler, globals)
		accProcessFindByAccNumberHandler := makeGlobalsHandler(bebber.AccProcessFindByAccNumberHandler, globals)
		readDocFileHandler := makeGlobalsHandler(bebber.ReadDocFileHandler, globals)
		docNumberProposalCurrHandler := makeGlobalsHandler(bebber.DocNumberProposalCurrHandler, globals)
		docNumberProposalChangeHandler := makeGlobalsHandler(bebber.DocNumberProposalChangeHandler, globals)
		docNumberProposalNextHandler := makeGlobalsHandler(bebber.DocNumberProposalNextHandler, globals)

		router := gin.Default()

		htmlDir := path.Join(config["PUBLIC_DIR"], "html")
		router.Use(bebber.Serve("/", bebber.LocalFile(htmlDir, false)))
		router.GET("/User/:name", authHandler, userHandler)
		router.POST("/Login", loginHandler)
		router.POST("/SearchDocs", authHandler, searchDocsHandler)
		router.POST("/Doc", authHandler, docMakeHandler)
		router.GET("/Doc/:name", authHandler, docReadHandler)
		router.PATCH("/Doc", authHandler, docChangeHandler)
		router.DELETE("/Doc/:name", authHandler, docRemoveHandler)
		router.PATCH("/DocRename", authHandler, docRenameHandler)
		router.PATCH("/DocLabels", authHandler, docAppendLabelsHandler)
		router.DELETE("/DocLabels/:name/:label", authHandler, docRemoveLabelsHandler)
		router.PATCH("/DocNumbers", authHandler, docAppendDocNumbersHandler)
		router.DELETE("/DocNumbers/:name/:number", authHandler, docRemoveDocNumberHandler)
		router.POST("/AccProcess", authHandler, accProcessMakeHandler)
		router.GET("/AccProcess/FindByDocNumber/:number", authHandler, accProcessFindByDocNumberHandler)
		router.GET("/AccProcess/FindByAccNumber/:from/:to/:number", authHandler, accProcessFindByAccNumberHandler)
		router.GET("/ReadDocFile/:name", authHandler, readDocFileHandler)
		router.GET("/DocNumberProposal", authHandler, docNumberProposalCurrHandler)
		router.GET("/DocNumberProposal/Next", authHandler, docNumberProposalNextHandler)
		router.PUT("/DocNumberProposal", authHandler, docNumberProposalChangeHandler)
		router.Static("/public", config["PUBLIC_DIR"])
		router.Static("/pdfviewer", config["PDFVIEWER_PUBLIC"])
	*/

	router := gin.New()

	// ----------------------
	//
	// THE BEGIN OF A NEW API
	//
	// ----------------------
	gumspecs.AppName = common.AppName
	httpServer := gumspecs.ReadHTTPServer()
	specs := common.LoadSpecs()
	gDB := common.InitMySQL()

	labels.AddTables(gDB)
	docs.AddTables(gDB)
	dbVars.AddTables(gDB)
	accountingData.AddTables(gDB)

	sessionStore := smem.NewStore()
	userStore := gumauth.SQLUserStore{gDB.Db}

	v1 := router.Group("/v1")

	signIn := gumauth.SignIn(sessionStore, userStore)
	v1.GET("/sign_in/:name/:password", signIn)

	g := v1.Group("/labels")
	{
		g.POST("/", gumwrap.Gorp(labels.CreateLabel, gDB))
		g.GET("/", gumwrap.Gorp(labels.ReadAllLabels, gDB)) // ?name="XX"
		g.GET("/:labelID", gumwrap.Gorp(labels.ReadOneLabel, gDB))
		g.DELETE("/:labelID", gumwrap.Gorp(labels.DeleteLabel, gDB))

		g.GET("/:labelID/docs", gumwrap.Gorp(docs.FindDocsWithLabelHandler, gDB))

	}

	g = v1.Group("/docs")
	{
		g.POST("/", gumwrap.Gorp(docs.CreateDocHandler, gDB))
		g.GET("/:docID", gumwrap.Gorp(docs.ReadOneDocHandler, gDB))
		g.PUT("/:docID", gumwrap.Gorp(docs.UpdateDocHandler, gDB))
		g.PATCH("/:docID/name", gumwrap.Gorp(docs.UpdateDocNameHandler, gDB))

		g.POST("/doc_numbers", gumwrap.Gorp(docs.CreateDocNumberHandler, gDB))
		g.GET("/:docID/doc_numbers", gumwrap.Gorp(docs.ReadAllDocNumbersHandler, gDB))
		g.DELETE("/:docID/doc_numbers/:docNumber", gumwrap.Gorp(docs.DeleteDocNumberHandler, gDB))

		g.POST("/account_data", gumwrap.Gorp(docs.CreateDocAccountDataHandler, gDB))
		g.GET("/:docID/account_data", gumwrap.Gorp(docs.ReadOneDocAccountDataHandler, gDB))
		g.PUT("/:docID/account_data", gumwrap.Gorp(docs.UpdateDocAccountDataHandler, gDB))

		g.POST("/labels", gumwrap.Gorp(docs.JoinLabelHandler, gDB))
		g.GET("/:docID/labels/", gumwrap.Gorp(docs.FindAllLabelsOfDocHandler, gDB))
		g.DELETE("/:docID/labels/:labelID", gumwrap.Gorp(docs.DetachLabelHandler, gDB))

		g.GET("/:docID/accounting_data", gumwrap.Gorp(docs.FindAllAccountingDataOfDocHandler, gDB))
	}

	g = v1.Group("/doc_number_proposals")
	{
		g.GET("/", gumwrap.Gorp(docNumberProposal.ReadDocNumberProposalHandler, gDB))
		g.GET("/next", gumwrap.Gorp(docNumberProposal.NextDocNumberProposalHandler, gDB))
		g.PUT("/", gumwrap.Gorp(docNumberProposal.UpdateDocNumberProposalHandler, gDB))
	}

	g = v1.Group("/accounting_data")
	{
		g.POST("/", gumwrap.Gorp(accountingData.CreateAccountingDataHandler, gDB))
	}

	g = v1.Group("/search")
	{
		g.POST("/docs", gumwrap.Gorp(docs.SearchDocsHandler, gDB))
	}

	router.GET("/docfile/:name", func(c *gin.Context) { docs.ReadDocFileHandler(c, specs) })
	htmlDir := path.Join(specs.Public, "html")
	router.Use(bebber.Serve("/", bebber.LocalFile(htmlDir, false)))
	router.Static("/public", specs.Public)
	router.Static("/pdfviewer", specs.PDFViewerPublic)

	router.Run(httpServer.String())

	//router.GET("/LoadBox/:boxname", bebber.Auth(), bebber.LoadBox)
	//router.POST("/AddTags", bebber.Auth(), bebber.AddTags)
	//router.GET("/LoadFile/:boxname/:filename", bebber.Auth(), bebber.LoadFile)
	//router.POST("/MoveFile", bebber.Auth(), bebber.MoveFile)
}

func StartAccServer(valid bool) {
	/*
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
	*/
}

func validCSV(valid bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("validCSV", valid)
		c.Next()
	}
}

func SetupDefault() (bebber.Config, bebber.MongoDBConn) {
	config := bebber.Config{}
	config["FILES_DIR"] = bebber.GetSettings("DOCMA_FILES")
	config["PUBLIC_DIR"] = bebber.GetSettings("DOCMA_PUBLIC")
	config["HTTP_HOST"] = bebber.GetSettings("DOCMA_HTTP_HOST")
	config["HTTP_PORT"] = bebber.GetSettings("DOCMA_HTTP_PORT")
	config["MONGODB_HOST"] = bebber.GetSettings("DOCMA_DB_SERVER")
	config["MONGODB_DBNAME"] = bebber.GetSettings("DOCMA_DB_NAME")
	config["PDFVIEWER_PUBLIC"] = bebber.GetSettings("DOCMA_PDFVIEWER_PUBLIC")

	dialInfo := &mgo.DialInfo{
		Addrs: []string{config["MONGODB_HOST"]},
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	conn := bebber.MongoDBConn{
		DialInfo: dialInfo,
		Session:  session,
		DBName:   config["MONGODB_DBNAME"],
	}

	return config, conn
}
