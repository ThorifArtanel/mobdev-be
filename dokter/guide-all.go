package dokter

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"mobdev.com/common"
)

type GuideAllReturn struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Logo    string `json:"logo"`
	Created string `json:"created"`
}

func GuideAll(c *gin.Context) {
	result := []GuideAllReturn{}

	db, err := common.DbConn()
	if err != nil {
		log.Print(err)
		log.Printf("Failed to create client %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed to create client : " + err.Error()})
		return
	}

	log.Print("Query #1 - All Guide")
	Q := `
		SELECT 
			guide_id,
			guide_title,
			guide_logo,
			created_dt
		FROM public.guide
	`
	rows, err := db.Query(Q)
	if err != nil {
		log.Print("db access error : " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Failed Querying Guide Data"})
		return
	}

	for rows.Next() {
		data := GuideAllReturn{}
		err = rows.Scan(&data.Id, &data.Title, &data.Logo, &data.Created)
		if err != nil {
			log.Print("db access error : " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "Failed Parsing Guide Data"})
			return
		}

		result = append(result, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "Success",
		"guide":  result,
	})
}
