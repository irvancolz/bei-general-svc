package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	BaseURL       string
	Authorization string
	ContentType   string
}

func (m *Request) Get(url string, target interface{}) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", m.BaseURL, url), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", m.ContentType)
	req.Header.Set("Authorization", m.Authorization)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(target)
}

func (m *Request) Post(url string, body, target interface{}) error {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", m.BaseURL, url), bytes.NewBuffer(bodyJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", m.ContentType)
	req.Header.Set("Authorization", m.Authorization)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(target)
}
