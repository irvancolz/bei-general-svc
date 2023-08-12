package utilities

import (
	"be-idx-tsg/internal/app/httprest/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	{
	}
	Request, err := http.NewRequest("GET", url, bodyReq)
	Request.Header.Add("authorization", token)
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
	body, err := io.ReadAll(resp.Body)

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

func GetParameterAdminImageExtension(c *gin.Context) (*APIResponseInterface, error) {
	err_host := godotenv.Load(".env")
	if err_host != nil {
		fmt.Println(err_host)
	}
	host := os.Getenv("SERVICE_AUTH_HOST")
	url := host + "/get-parameter-admin-by-key?key=format_file_logo"
	tokens, _ := c.Get("token")
	var payload string
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("[AQI] [err] [GetRequest][Payload] ", err)
	}
	bodyReq := bytes.NewReader(payloadBytes)
	token := tokens.(string)
	Request, err := http.NewRequest("GET", url, bodyReq)
	Request.Header.Add("authorization", token)
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
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("[AQI] [err] [GetRequest][ReadAll]", err)
		return nil, err
	}

	datas := &APIResponseInterface{}
	errorUM := json.Unmarshal([]byte(body), datas)
	if errorUM != nil {
		log.Println("[AQI] [err] [GetRequest][errorUM]", errorUM)
	}

	return datas, errorUM
}

type UserDetailResponse struct {
	Code    int64      `json:"code"`
	Message string     `json:"message"`
	Data    model.User `json:"data"`
}

func GetUserNameByID(c *gin.Context, id string) string {
	err_host := godotenv.Load(".env")
	if err_host != nil {
		fmt.Println(err_host)
	}
	host := os.Getenv("SERVICE_AUTH_HOST")

	url := host + "/management-user-get-user-by-id?id=" + id
	tokens, _ := c.Get("token")

	Request, err := http.NewRequest("GET", url, nil)

	Request.Header.Add("authorization", tokens.(string))
	Request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("failed to create request to auth : ", err)
		return ""
	}

	resp, err := http.DefaultClient.Do(Request)
	if err != nil {
		log.Println("failed to get user response from auth: ", err)
		return ""
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Println("failed to get user detail, an error occured when try to get data")
		return ""
	}

	if err != nil {
		log.Println("failed to read data user results : ", err)
		return ""
	}

	result := UserDetailResponse{}
	err_marshall := json.Unmarshal(data, &result)
	if err_marshall != nil {
		log.Println("failed to convert users to expected struct :", err_marshall)
		return ""
	}

	return result.Data.Username
}

type UserRolesResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func GetUserRoles(c *gin.Context, id string) string {
	err_host := godotenv.Load(".env")
	if err_host != nil {
		fmt.Println(err_host)
	}
	host := os.Getenv("SERVICE_AUTH_HOST")

	url := host + "/get-user-roles-by-id?id=" + id
	tokens, _ := c.Get("token")

	Request, err := http.NewRequest("GET", url, nil)

	Request.Header.Add("authorization", tokens.(string))
	Request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("failed to create request to auth : ", err)
		return ""
	}

	resp, err := http.DefaultClient.Do(Request)
	if err != nil {
		log.Println("failed to get user response from auth: ", err)
		return ""
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Println("failed to get user detail, an error occured when try to get data")
		return ""
	}

	if err != nil {
		log.Println("failed to read data user results : ", err)
		return ""
	}

	result := UserRolesResponse{}
	err_marshall := json.Unmarshal(data, &result)
	if err_marshall != nil {
		log.Println("failed to convert users to expected struct :", err_marshall)
		return ""
	}

	return result.Data
}
