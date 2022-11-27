package dokter

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"mobdev.com/common"
)

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func DokterAuth(c *gin.Context) {
	db, err := common.DbConn()
	if err != nil {
		log.Print(err)
		log.Printf("Failed to create client %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed to create client : " + err.Error()})

		return
	}

	log.Print("Parsing Payload")
	var req Auth
	if err := c.BindJSON(&req); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "Failed Parsing Payload",
		})
		return
	}

	log.Println("Query #1 AuthLogin")
	var pw string
	var raw_token common.Token
	Q := `
		SELECT 
			dkt.dkt_id,
			dkt.dkt_name,
			dkt.dkt_password
		FROM public.user_dokter dkt
		WHERE dkt_email=$1;
	`
	err = db.QueryRow(Q, req.Email).Scan(&raw_token.Id, &raw_token.UserName, &pw)
	raw_token.UserGroup = "Dokter"
	switch {
	case err == sql.ErrNoRows:
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Incorrect Username or Password"})

			return
		}
	case err != nil:
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "Failed Retrieving User Data"})

			return
		}
	}
	if !common.ComparePasswords(pw, []byte(req.Password)) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "Incorrect Username or Password"})

		return
	}

	token, err := common.CreateToken(raw_token)
	if err != nil {
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "Failed Creating Token"})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Login Approved",
		"token":  token,
	})
}
