package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func UpdateFormAttachmentFileStatus(c *gin.Context, filename string) {
	err_host := godotenv.Load(".env")
	if err_host != nil {
		fmt.Println(err_host)
	}
	host := os.Getenv("SERVICE_FORM_HOST")

	var bodyReq struct {
		FileName string `json:"file_name"`
	}

	bodyReq.FileName = filename
	bodyReqByte, errMarshall := json.Marshal(bodyReq)
	if errMarshall != nil {
		log.Println("failed convert request body to json :", errMarshall)
	}
	bodyReqReader := bytes.NewReader(bodyReqByte)

	url := host + "/form-value-attachment/delete"
	tokens, _ := c.Get("token")

	Request, err := http.NewRequest("DELETE", url, bodyReqReader)

	Request.Header.Add("authorization", tokens.(string))
	Request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("failed to create request to form : ", err)
		return
	}

	resp, err := http.DefaultClient.Do(Request)
	if err != nil {
		log.Println("failed to delete attachment in form services: ", err)
		return
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Println("failed toto delete attachment, an error occured when try to delete data")
		return
	}

	if err != nil {
		log.Println("failed to read data user results : ", err)
		return
	}
	log.Println(string(data))
}
