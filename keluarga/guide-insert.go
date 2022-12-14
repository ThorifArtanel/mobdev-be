package keluarga

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/jackc/pgconn"
	"mobdev.com/common"
)

type KeluargaGuideInsertReq struct {
	Title   string `json:"title" validate:"required"`
	Subject string `json:"subject" validate:"required,AlphaSpace"`
	Content string `json:"content" validate:"required"`
}

func KeluargaGuideInsert(c *gin.Context) {
	log.Print("Parsing Payload")
	var req KeluargaGuideInsertReq
	if err := c.Bind(&req); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "Failed Parsing Payload"})
		return
	}

	validate := validator.New()
	validate.RegisterValidation("AlphaSpace", common.ValidateStringWhitespace)
	err := validate.Struct(req)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"msg":    "Invalid Parameters " + err.Error()})
		return
	}

	db, err := common.DbConn()
	if err != nil {
		log.Print(err)
		log.Printf("Failed to create client %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed to create client : " + err.Error()})
		return
	}

	log.Print("Query #1 - Detail Guide")
	Q := `
		INSERT INTO public.guide(
			guide_id,
			guide_title,
			guide_subject,
			guide_desc,
			guide_content,
			guide_logo,
			created_by,
			updated_by,
			dkt_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $7, $7
		);
	`
	desc := common.TruncateToLength(req.Content, common.GetDescTruncLen())
	rndm := rand.Intn(999)
	logo := "id/" + strconv.Itoa(rndm) + "/200/200"
	_, err = db.Exec(Q, common.GenerateUUID(), req.Title, req.Subject, desc, req.Content, logo, "ADMIN")
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{
					"code": http.StatusConflict,
					"msg":  "Conflict Of Guide Kode"})
				return
			}
		}

		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed Inserting Guide Data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
	})
}
