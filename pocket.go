package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/viper"
	"net/http"
)

// Set an init function which will get the config via viper and set the Cobra
// instance going

// TODO:
// 1. Update the Yaml file with the latest_post value
// 2. Set up commands ( ie. add, remove, open, login etc )
// 3. URL encode (escape) urls sent to pocket API
// 4. Write documentation

type feedItem struct {
	URL        string `mapstructure:"url"`
	Tags       string `mapstructure:"tags"`
	LatestPost string `mapstructure:"latest_post"`
}

type feeds struct {
	Feeds []feedItem
}

func main() {
	// Set the name for the config file
	viper.SetConfigName("pocket")

	// Set a path for the config to be found in
	viper.AddConfigPath(".")

	// Set another path
	viper.AddConfigPath("$HOME/.config/pocket")

	// Read what the config says
	err := viper.ReadInConfig()

	// Check for errors
	if err != nil {
		// Print the error to the console
		fmt.Println(err)
	}

	// Create a new struct for feeds
	f := feeds{}

	// Populate this structure with data from the yaml
	err = viper.Unmarshal(&f)

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
				send(element.Link, tags)
			}
		}
		// Add a space for prettyness
		fmt.Println()

		element.LatestPost = feed.Items[0].GUID
		// viper.Set("feeds.latestpost", "test")
		// err = viper.WriteConfig()
	}

}

func send(URL string, Tags string) {
	// PocketConfig : The Pocket API will need a key and access token to work with
	type PocketConfig struct {
		ConsumerKey string `mapstructure:"consumer_key"`
		AccessToken string `mapstructure:"access_token"`
	}

	// PostRequest : To construct some json which can be sent to the Pocket API
	type PostRequest struct {
		URL   string `json:"url"`
		Key   string `json:"consumer_key"`
		Token string `json:"access_token"`
		Tags  string `json:"tags"`
	}

	// Get the pocket specific data
	Config := PocketConfig{}

	// err = yaml.Unmarshal(yamlFile, &P)
	err := viper.Unmarshal(&Config)

	// Check for errors
	if err != nil {
		// Print the error to the console
		fmt.Println(err)
	}

	fmt.Println(Config.ConsumerKey)

	// Pocket API enpoint for adding items
	url := "https://getpocket.com/v3/add"

	// Gather data use use in POST request
	jsonStr := &PostRequest{URL, Config.ConsumerKey, Config.AccessToken, Tags}

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
