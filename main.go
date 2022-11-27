package main

import (
	"github.com/gin-gonic/gin"
	"mobdev.com/common"
	dkt "mobdev.com/dokter"
)

func main() {
	r := gin.Default()
	r.Use(common.CORSMiddleware())

	// For Endpoints That Doesn't Need Authorization
	r.POST("/dokter/auth", dkt.DokterAuth)

	// For Endpoints That Need Authorization
	authorized := r.Group("/")
	{
		authorized.Use(common.AuthMiddleware())
		// authorized.GET("/cms/user", UserAll)
	}

	r.Run(":80")
}
