package main

import (
	"github.com/gin-gonic/gin"
	"mobdev.com/common"

	dkt "mobdev.com/dokter"
	klg "mobdev.com/keluarga"
)

func main() {
	r := gin.Default()
	r.Use(common.CORSMiddleware())

	// Doctor Endpoints
	r.POST("/dokter/auth", dkt.DokterAuth)

	r.GET("/dokter/:id", dkt.DokterOne)
	r.PUT("/dokter/:id", dkt.DokterUpdate)
	r.POST("/dokter", dkt.DokterInsert)

	r.GET("/dokter/guide", dkt.DokterGuideAll)
	r.POST("/dokter/guide", dkt.DokterGuideInsert)

	// Keluarga Endpoints
	r.POST("/dokter/auth", klg.KeluargaAuth)

	r.GET("/keluarga/guide", klg.KeluargaGuideAll)
	r.POST("/keluarga/guide", klg.KeluargaGuideInsert)

	r.Run(":80")
}
