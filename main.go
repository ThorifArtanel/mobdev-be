package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"mobdev.com/common"

	dkt "mobdev.com/dokter"
)

func main() {
	log.Print(common.GetDBURL())
	r := gin.Default()
	r.Use(common.CORSMiddleware())

	// For Endpoints That Doesn't Need Authorization
	r.POST("/dokter/auth", dkt.DokterAuth)

	r.GET("/dokter/:id", dkt.DokterOne)
	r.PUT("/dokter/:id", dkt.DokterUpdate)
	r.POST("/dokter", dkt.DokterInsert)

	r.GET("/dokter/guide", dkt.DokterGuideAll)
	r.POST("/dokter/guide", dkt.DokterGuideInsert)

	// For Endpoints That Need Authorization
	authorized := r.Group("/")
	{
		authorized.Use(common.AuthMiddleware())
		// authorized.GET("/cms/user", UserAll)
	}

	r.Run(":80")
}
