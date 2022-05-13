package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	    "io"
	    "log"
	"net/http"
)

type BagConfig struct {
    BaseUrl string `mapstructure:"baseUrl"`
    AccessToken string `mapstructure:"access_token"`
    ClientID string `mapstructure:"client_id"`
    ClientSecret string `mapstructure:"client_secret"`
    UserName string `mapstructure:"username"`
    Password string `mapstructure:"password"`
}

func auth() string {
	type PostRequest struct {
        GrantType string `json:"grant_type"`
        ClientID string `json:"client_id"`
        ClientSecret string `json:"client_secret"`
        UserName string `json:"username"`
        Password string `json:"password"`
	}


	Config := BagConfig{}
	err := viper.Unmarshal(&Config)

	// Check for errors
	if err != nil {
		// Print the error to the console
		panic(err)
	}
    bearer := "Bearer " + Config.AccessToken

	// Pocket API enpoint for adding items
	baseURL := Config.BaseUrl + "api/entries.json"

	// Configure a new request using the URL and Json
	req, err := http.NewRequest("GET", baseURL, nil)
	req.Header.Set("Content-Type", "application/json")
    req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)

    // Check for errors
    if err != nil {
        // Print the error to the console
        panic(err)
    }

    // fmt.Println("🍺  this is run ", resp.Status)
    // fmt.Println("🍺  this is run ", resp)
	// Close the response
	defer resp.Body.Close()

    if resp.StatusCode != 200 {
        // fmt.Println("🍺  Do Auth ", resp.StatusCode)
        // Pocket API enpoint for adding items
        baseURL := Config.BaseUrl + "oauth/v2/token"

        jsonStr := &PostRequest{
            "password",
            Config.ClientID,
            Config.ClientSecret,
            Config.UserName,
            Config.Password,
        }

        // fmt.Println("🍺  jsonStr ", jsonStr)

        // Json encode this data
        b, err := json.Marshal(jsonStr)

        // fmt.Println("🍺  jsonStr ", bytes.NewBuffer(b))
        // Configure a new request using the URL and Json
        req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(b))
        req.Header.Add("Content-Type", "application/json")
        req.Header.Add("Accept", "application/json, */*;q=0.5")
        req.Header.Add("Authorization", bearer)

        // debug(httputil.DumpRequestOut(req, true))

        client := &http.Client{}
        resp, err := client.Do(req)
        fmt.Println("🍺  after auth ", resp.Body)

        // Check for errors
        if err != nil {
            // Print the error to the console
            panic(err)
        }
        // Close the response
        defer resp.Body.Close()


        fmt.Println("🍺  The Response ", resp.Body)

        b, err = io.ReadAll(resp.Body)
        if err != nil {
            log.Fatalln(err)
        }

        return string(b)
    }
    return ""
}
