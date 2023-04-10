package helper

import (
	"be-idx-tsg/internal/app/httprest/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"os"

	"github.com/joho/godotenv"
)

func CheckPermission(token string, currentPermission *string) (*model.GetAuthResponses, error) {

	type payload_struct struct {
		Module *string `json:"module"`
	}

	payload := payload_struct{
		Module: currentPermission,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	err_host := godotenv.Load(".env")
	if err_host != nil {
		fmt.Println(err_host)
	}
	host := os.Getenv("SERVICE_MIDDLEWARE_HOST")

	bodyReq := bytes.NewReader(payloadBytes)
	url := host + "/check-permission"
	// tokens, _ := c.Get("token")
	log.Println(url)
	Request, err := http.NewRequest("POST", url, bodyReq)
	Request.Header.Add("authorization", token)
	Request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("error")
		log.Println(err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(Request)
	if err != nil {
		log.Println("[TAP-debug] [err] [PostRequest][Do]", err)
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("[TAP-debug] [err] [PostRequest][ReadAll]", err)
		log.Println(err)
		return nil, err
	}

	datas := &model.GetAuthResponses{}
	errorUM := json.Unmarshal([]byte(body), datas)

	return datas, errorUM
}
