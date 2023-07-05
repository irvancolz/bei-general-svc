package utilities

import (
	"be-idx-tsg/internal/app/httprest/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type ContactpersonCompaniesSyncRequest struct {
	Code    int64                                      `json:"code"`
	Message string                                     `json:"message"`
	Data    []model.ContactPersonSyncCompaniesResource `json:"data"`
}

func GetLatestABCompanies(c *gin.Context) ([]model.ContactPersonSyncCompaniesResource, error) {
	err_host := godotenv.Load(".env")
	if err_host != nil {
		fmt.Println(err_host)
	}
	host := os.Getenv("SERVICE_AB_HOST")

	url := host + "/get-ab-to-sync-with-contact-person"
	tokens, _ := c.Get("token")

	Request, err := http.NewRequest("GET", url, nil)

	Request.Header.Add("authorization", tokens.(string))
	Request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("failed to get data from : ", url)
		log.Println("failed to get user email : ", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(Request)
	if err != nil {
		log.Println("failed to get user response: ", err)
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Println("failed to get anggota bursa companies from ab list, error occured when try to get data")
		return nil, errors.New("failed to get anggota bursa companies from ab list, error occured when try to get data")
	}

	if err != nil {
		log.Println("failed to read data user results : ", err)
		return nil, err
	}

	result := ContactpersonCompaniesSyncRequest{}
	err_marshall := json.Unmarshal(data, &result)
	if err_marshall != nil {
		return nil, err_marshall
	}

	return result.Data, nil
}

func GetLatestParticipantCompanies(c *gin.Context) ([]model.ContactPersonSyncCompaniesResource, error) {
	err_host := godotenv.Load(".env")
	if err_host != nil {
		fmt.Println(err_host)
	}
	host := os.Getenv("SERVICE_PARTICIPANT_HOST")

	url := host + "/get-participant-to-sync-with-contact-person"
	tokens, _ := c.Get("token")

	Request, err := http.NewRequest("GET", url, nil)

	Request.Header.Add("authorization", tokens.(string))
	Request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("failed to get data from : ", url)
		log.Println("failed to get user email : ", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(Request)
	if err != nil {
		log.Println("failed to get user response: ", err)
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Println("failed to get anggota bursa companies from ab list, error occured when try to get data")
		return nil, errors.New("failed to get anggota bursa companies from ab list, error occured when try to get data")
	}

	if err != nil {
		log.Println("failed to read data user results : ", err)
		return nil, err
	}

	result := ContactpersonCompaniesSyncRequest{}
	err_marshall := json.Unmarshal(data, &result)
	if err_marshall != nil {
		return nil, err_marshall
	}

	return result.Data, nil
}
