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

type DokterInsertReq struct {
	Email       string `json:"email" validate:"required,email"`
	Name        string `json:"name"  validate:"required,AlphaSpace"`
	NIP         string `json:"nip" validate:"omitempty,int,min=1, max=999999999999999999"`
	Phone       string `json:"phone" validate:"required,int"`
	NewPass     string `json:"new_pass" validate:"required,eqfield=ConfirmPass"`
	ConfirmPass string `json:"confirm_pass" validate:"required,eqfield=NewPass"`
}

func DokterInsert(c *gin.Context) {
	token := c.MustGet("Token").(*common.Token)

	if token.UserGroup != common.DokterRole() {
		log.Print("invalid role")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "Invalid Role"})
		return
	}

	log.Print("Parsing Payload")
	var req DokterInsertReq
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

	strong := common.IsStrongPassword(req.NewPass)
	if !strong {
		log.Print(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"msg":    "Password minimal terdiri dari 8 karakter dengan satu huruf besar, angka, dan karakter.",
		})
		return
	}
	pass := common.HashAndSalt([]byte(req.NewPass))

	log.Print("Query #1 - Detail Dokter")
	Q := `
		INSERT INTO public.user_dokter(
			dkt_id,
			dkt_email,
			dkt_name,
			dkt_nip,
			dkt_phone,
			dkt_password,
			created_by,
			updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $7
		);
	`
	_, err = db.Exec(Q, common.GenerateUUID(), req.Email, req.Name, req.NIP, req.Phone, pass, token.Id)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{
					"code": http.StatusConflict,
					"msg":  "Conflict Of Dokter Kode"})
				return
			}
		}

		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed Inserting Dokter Data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
	})
}
