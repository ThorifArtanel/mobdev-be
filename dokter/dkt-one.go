package dokter

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"mobdev.com/common"
)

type DokterOneReturn struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	NIP     string `json:"nip"`
	Phone   string `json:"phone"`
	Picture string `json:"picture"`
}

func DokterOne(c *gin.Context) {
	result := DokterOneReturn{}

	id := c.Param("id")

	db, err := common.DbConn()
	if err != nil {
		log.Print(err)
		log.Printf("Failed to create client %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed to create client : " + err.Error()})
		return
	}

	log.Print("Query #1 - Detail Dokter")
	Q := `
		SELECT 
			dkt_email,
			dkt_name,
			dkt_nip,
			dkt_phone,
			dkt_profile_pic
		FROM public.user_dokter
		WHERE dkt_id=$1;
	
	`
	err = db.QueryRow(Q, id).Scan(&result.Email, &result.Name, &result.NIP, &result.Phone, &result.Picture)
	result.Id = id
	switch {
	case err == sql.ErrNoRows:
		log.Print("dokter data not found")
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "Dokter Data Not Found"})
		return
	case err != nil:
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed Retrieving Dokter Data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
		"dokter": result,
	})
}
