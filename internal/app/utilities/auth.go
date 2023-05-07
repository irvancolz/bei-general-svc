package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)


func GetAllRole(c *gin.Context) (*APIResponse, error) {
	err_host := godotenv.Load(".env")
	if err_host != nil {
		fmt.Println(err_host)
	}
	host := os.Getenv("SERVICE_AUTH_HOST")
	url := host + "/get-all-user-role"
	tokens, _ := c.Get("token")
	var payload string
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("[AQI] [err] [GetRequest][Payload] ", err)
	}
	bodyReq := bytes.NewReader(payloadBytes)
	token := tokens.(string)
	{}
	Request, err := http.NewRequest("GET", url,bodyReq)	
	Request.Header.Add("authorization",token )
	Request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("[AQI] [err] [GetRequest][Wraps] ", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(Request)
	if err != nil {
		log.Println("[AQI] [err] [GetRequest][Do]", err)
		
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("[AQI] [err] [GetRequest][ReadAll]", err)
		return nil, err
	}

	datas := &APIResponse{}
	errorUM := json.Unmarshal([]byte(body), datas)
	if errorUM != nil {
		log.Println("[AQI] [err] [GetRequest][errorUM]", errorUM)
	}

	return datas, errorUM
}
