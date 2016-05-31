package main

import (
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
)

func main() {
	StartDefaultServer()
}

func StartDefaultServer() {

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
		g.PATCH("/:docID/name", func(c *gin.Context) { docs.UpdateDocNameHandler(c, gDB, specs.Files) })
		g.DELETE("/:docID", gumwrap.Gorp(docs.RemoveDocHandler, gDB))

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
