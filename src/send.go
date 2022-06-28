package gobag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
)

func Send(postURL string, Tags string, AccessToken string) {

	// PostRequest : To construct some json which can be sent to the Pocket API
	type PostRequest struct {
		URL   string `json:"url"`
		Tags  string `json:"tags"`
	}

	// Get the pocket specific data
	Config := BagConfig{}

	// Populate struct with config data
	err := viper.Unmarshal(&Config)

	// Check for errors
	if err != nil {
		// Print the error to the console
		panic(err)
	}

	// Pocket API enpoint for adding items
	baseURL := Config.BaseUrl + "api/entries.json"

	// Encode string as a URL
	// https://getpocket.com/developer/docs/v3/add - Best Practices
	u, err := url.Parse(postURL)
	// Convert encoded URL back to a string
	URL := u.String()

	// Gather data use use in POST request
	jsonStr := &PostRequest{URL, Tags}

	// Json encode this data
	b, err := json.Marshal(jsonStr)

    bearer := "Bearer " + AccessToken

	// Configure a new request using the URL and Json
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
    req.Header.Add("Authorization", bearer)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)

	// fmt.Println("🐛  ", req)
	// fmt.Println("🐛  ", resp)

	// Handle errors
	if err != nil {
		panic(err)
	}

	// Close the response
	defer resp.Body.Close()

	// Inform the user which URL has been sent
	fmt.Println("🚀  ", URL)
}
