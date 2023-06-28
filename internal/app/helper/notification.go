package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type NotificationData struct {
	Title string `json:"title"`
	Date  string `json:"date"`
}

type CreateSingleNotificationParam struct {
	UserID string           `json:"user_id"`
	Data   NotificationData `json:"data"`
	Link   string           `json:"link"`
	Type   string           `json:"type"`
}

func CreateSingleNotification(c *gin.Context, param CreateSingleNotificationParam) {
	data, _ := json.Marshal(param)

	token, _ := c.Get("token")

	url := os.Getenv("SERVICE_AUTH_HOST") + "/create-new-notifications"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))

	req.Header.Add("authorization", token.(string))

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println("failed to create notifications", err)
		return
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("[TAP-debug] [err] [PostRequest][Do]", err)
		return
	}

	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		log.Println("[TAP-debug] [err] [PostRequest][ReadAll]", err)
		log.Println(body)
		return
	}
}

func GetCurrentTime() string {
	t, _ := TimeIn(time.Now(), "Asia/Jakarta")

	return t.Format("2006-01-02 15:04:05")
}
