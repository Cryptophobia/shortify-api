package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type APIResponse struct {
	StatusMessage string `json:statusmessage`
	ShortURL      string `json:shorturl`
}

const urlCreate = "http://localhost:5000/Create"

// Testing the Create route
func TestShortifyCreate(t *testing.T) {
	payloadCreate := []byte(`{"longurl": "https://github.com/teamhephy/workflow/"}`)

	req, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(payloadCreate))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check if response is OK Status
	assert.Equal(t, "200 OK", resp.Status)
	// Check if the expected shortString is 8 characters long.
	body, _ := ioutil.ReadAll(resp.Body)
	responseObj := &APIResponse{}
	json.Unmarshal([]byte(body), responseObj)

	assert.Equal(t, 8, len(responseObj.ShortURL))
	assert.Equal(t, "OK", responseObj.StatusMessage)
}

// Testing both Create and Get routes
func TestSHortifyCreateGet(t *testing.T) {
	payloadCreate := []byte(`{"longurl": "https://teamhephy.com"}`)

	req, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(payloadCreate))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check if response is OK Status
	assert.Equal(t, "200 OK", resp.Status)
	// Check if the expected shortString is 8 characters long.
	body, _ := ioutil.ReadAll(resp.Body)
	responseObj := &APIResponse{}
	json.Unmarshal([]byte(body), responseObj)

	assert.Equal(t, 8, len(responseObj.ShortURL))
	assert.Equal(t, "OK", responseObj.StatusMessage)

	// Use the returned object to make a new get request
	urlGet := fmt.Sprintf("%v%v", "http://localhost:5000/", responseObj.ShortURL)

	req, err = http.NewRequest("GET", urlGet, nil)
	if err != nil {
		t.Fatal(err)
	}

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, "200 OK", resp.Status)
	assert.Equal(t, true, strings.Contains(string(body), "<meta name=\"description\" content=\"Hephy Workflow puts a cool breeze in the Kubernetes sails!\">"))
}

// Testing the validation
func TestShortifyURLValidation(t *testing.T) {
	failPayloads := make(map[int][]byte)
	failPayloads[0] = []byte(`{"longurl": "https//github.com/teamhephy/workflow/"}`)
	failPayloads[1] = []byte(`{"longurl": "http:::/not.valid/a//a??a?b=&&c#hi"}`)
	failPayloads[2] = []byte(`{"longurl": "cloudflare.com"}`)
	failPayloads[3] = []byte(`{"longurl": "/foo/bar"}`)
	failPayloads[4] = []byte(`{"longurl": "http://"}`)

	for _, payload := range failPayloads {
		req, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		// Check if response is 400 Bad Request Status
		assert.Equal(t, "400 Bad Request", resp.Status)
		time.Sleep(20 * time.Millisecond)
	}
}

// Testing the validation
func TestShortifyURLVadationShort(t *testing.T) {
	failShortURLs := make(map[int]string)
	failShortURLs[0] = "8sa3efd"
	failShortURLs[1] = "0122"
	failShortURLs[2] = "¢¥®µÐ®µÐ"

	for _, shortURL := range failShortURLs {
		// Use the returned object to make a new get request
		urlGet := fmt.Sprintf("%v%v", "http://localhost:5000/", shortURL)
		req, err := http.NewRequest("GET", urlGet, strings.NewReader(shortURL))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		// Check if response is 400 Bad Request Status
		assert.Equal(t, "400 Bad Request", resp.Status)
		time.Sleep(20 * time.Millisecond)
	}
}
