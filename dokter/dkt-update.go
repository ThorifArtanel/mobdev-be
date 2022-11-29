package dokter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/jackc/pgconn"
	"mobdev.com/common"
)

type DokterUpdateReq struct {
	Name  string `json:"name"  validate:"required,AlphaSpace"`
	NIP   string `json:"nip" validate:"omitempty,int,min=1, max=999999999999999999"`
	Phone string `json:"phone" validate:"required,int"`
}

func DokterUpdate(c *gin.Context) {
	token := c.MustGet("Token").(*common.Token)
	id := c.Param("id")

	log.Print("Parsing Payload")
	var req DokterUpdateReq
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

	log.Print("Query #1 - Update Dokter")
	Q := `
		UPDATE public.user_dokter SET 
			dkt_name=$2,
			dkt_nip=$3,
			dkt_phone=$4,
			updated_by=$5
		WHERE dkt_id=$1;
	`
	_, err = db.Exec(Q, id, req.Name, req.NIP, req.Phone, token.Id)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{
					"code": http.StatusConflict,
					"msg":  "Conflict Of Dokter Id"})
				return
			}
		}

		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed Updating Dokter Data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
	})
}
