package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Result struct {
	Reference string `json:"reference"`
	Endpoint  string `json:"endpoint"`
	Response  string `json:"response"`
	Error     string `json:"error,omitempty"`
}

func CallAPIWithTimeout(url string, vrm string, contraventionDate string, timeout time.Duration, resultChan chan<- Result) {
	searchBody := map[string]string{
		"vrm":                vrm,
		"contravention_date": contraventionDate,
	}

	jsonData, err := json.Marshal(searchBody)
	if err != nil {
		resultChan <- Result{Reference: "", Endpoint: url, Error: err.Error()}
		return
	}

	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		resultChan <- Result{Reference: "", Endpoint: url, Error: err.Error()}
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		resultChan <- Result{Reference: "", Endpoint: url, Error: err.Error()}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultChan <- Result{Reference: "", Endpoint: url, Error: err.Error()}
		return
	}

	resultChan <- Result{
		Reference: uuid.New().String(),
		Endpoint:  url,
		Response:  string(body),
	}
}
