package utilities

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type GetDataFeedResponse struct {
	ID            string     `json:"id"`
	TypeID        string     `json:"type_id"`
	Text          string     `json:"text"`
	RelationId    *string    `json:"relation_id"`
	Data          *string    `json:"data"`
	UpdatedByName string     `json:"updated_by_name"`
	DateTime      *string    `json:"date_time"`
	Status        bool       `json:"status"`
	CreatedAt     *time.Time `json:"created_at"`
	CreatedBy     *string    `json:"created_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedBy     *string    `json:"updated_by"`
}

func GetFollowUpPersonel(c *gin.Context) []string {
	var result []string
	host := os.Getenv("SERVICE_DATAFEED_HOST")

	url := host + "/get-data-feed-by-type?type=Personel"
	tokens, _ := c.Get("token")

	Request, err := http.NewRequest("GET", url, nil)

	Request.Header.Add("authorization", tokens.(string))
	Request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("failed to create request to auth : ", err)
		return result
	}

	resp, err := http.DefaultClient.Do(Request)
	if err != nil {
		log.Println("failed to get user response from auth: ", err)
		return result
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Println("failed to get user detail, an error occured when try to get data")
		return result
	}

	if err != nil {
		log.Println("failed to read data user results : ", err)
		return result
	}

	var apiRes HTTPResp[[]GetDataFeedResponse]

	err_marshall := json.Unmarshal(data, &apiRes)
	if err_marshall != nil {
		log.Println("failed to convert users to expected struct :", err_marshall)
		return result
	}

	for _, res := range apiRes.Data {
		result = append(result, res.Text)
	}

	return result
}
