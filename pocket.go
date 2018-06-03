package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

// TODO:
// 1. Update the Yaml file with the latest_post value
// 2. Set up commands ( ie. add, remove, open, login etc )
// 3. URL encode (escape) urls sent to pocket API
// 4. Write documentation
// 5. (BONUS) use https://github.com/spf13/viper for the config

// FeedItem : Map each of the feeds data from yaml to Go
type FeedItem struct {
	URL          string `yaml:"url"`
	Tags         string `yaml:"tags"`
	LastModified string `yaml:"last_modified"`
	LatestPost   string `yaml:"latest_post"`
}

// Get all the feeds
type feeds struct {
	Feeds []FeedItem
}

func main() {
	// Read the yaml file
	yamlFile, err := ioutil.ReadFile("pocket.yaml")

	// Check for errors
	if err != nil {
		// Print the error to the console
		fmt.Println(err)
	}

	// Create a new struct for feeds
	f := feeds{}

	// Populate this structure with data from the yaml
	err = yaml.Unmarshal(yamlFile, &f)

	// Check for errors
	if err != nil {
		// Print the error to the console
		fmt.Println(err)
	}

	// Create a new RSS parser instance
	fp := gofeed.NewParser()

	// For wach of the feeds from the Yaml
	for _, element := range f.Feeds {

		// Indicate something is happeing
		fmt.Println("🍺  downloading ", element.URL)

		// parse the RSS url for that feed
		feed, _ := fp.ParseURL(element.URL)

		// print when that feed was last updated
		// fmt.Println("⏰  Last updated on ", feed.Updated)

		// Save the ID for the last post fetched
		latest := element.LatestPost

		// Save the tags from the yaml
		tags := element.Tags

		// For each of the items in that feed
		for _, element := range feed.Items {

			// If the current ID is the last one we got
			if element.GUID == latest {
				// Stop looking through the posts
				break
			} else {
				// Send all the newer posts to pocket
				sendToPocket(element.Link, tags)
			}
		}
		// Add a space for prettyness
		fmt.Println()
	}

}

// PocketConfig : The Pocket API will need a key and access token to work with
type PocketConfig struct {
	Key   string `yaml:"consumer_key"`
	Token string `yaml:"access_token"`
}

// PostRequest : To construct some json which can be sent to the Pocket API
type PostRequest struct {
	URL   string `json:"url"`
	Key   string `json:"consumer_key"`
	Token string `json:"access_token"`
	Tags  string `json:"tags"`
}

func sendToPocket(URL string, Tags string) {
	// Read the Yaml Configuration
	yamlFile, err := ioutil.ReadFile("pocket.yaml")

	// Get the pocket specific data
	P := PocketConfig{}
	err = yaml.Unmarshal(yamlFile, &P)

	// Check for errors
	if err != nil {
		// Print the error to the console
		fmt.Println(err)
	}

	// Pocket API enpoint for adding items
	url := "https://getpocket.com/v3/add"

	// Gather data use use in POST request
	jsonStr := &PostRequest{URL, P.Key, P.Token, Tags}

	// Json encode this data
	b, err := json.Marshal(jsonStr)

	// Configure a new request using the URL and Json
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)

	// Handle errors
	if err != nil {
		panic(err)
	}

	// Close the response
	defer resp.Body.Close()

	// Inform the user which URL has been sent
	fmt.Println("🚀  ", URL)
}
