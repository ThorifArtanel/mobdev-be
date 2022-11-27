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
	r.POST("/dokter/guide", dkt.DokterGuide)

	// For Endpoints That Need Authorization
	authorized := r.Group("/")
	{
		authorized.Use(common.AuthMiddleware())
		// authorized.GET("/cms/user", UserAll)
	}

	r.Run(":80")
}
