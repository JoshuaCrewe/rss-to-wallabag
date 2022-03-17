package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
)

// Set an init function which will get the config via viper and set the Cobra
// instance going
// https://github.com/spf13/cobra

// TODO:
// [X] 1. Update the Yaml file with the latest_post value
// [ ] 2. Set up commands ( ie. add, remove, open, login etc ) - Cobra
// [X] 3. URL encode (escape) urls sent to wallabag API
// [ ] 4. Write documentation

type feedItem struct {
	URL        string `mapstructure:"url"`
	Tags       string `mapstructure:"tags"`
	LatestPost string `mapstructure:"latestpost"`
}

type feeds struct {
	Feeds []feedItem
}

func main() {
	// Set the name for the config file
	viper.SetConfigName("bag")

	// Set a path for the config to be found in
	viper.AddConfigPath(".")

	// Set another path
	viper.AddConfigPath("$HOME/.config/bag")

	// Read what the config says
	err := viper.ReadInConfig()

	// Check for errors
	if err != nil {
		// Print the error to the console
		panic(err)
	}

	// Create a new struct for feeds
	f := feeds{}

	// Populate this structure with data from the yaml
	err = viper.Unmarshal(&f)

	// Check for errors
	if err != nil {
		// Print the error to the console
		panic(err)
	}

	// Create a new RSS parser instance
	fp := gofeed.NewParser()

	// Keep a count for the loops
	i := 0

	// For each of the feeds from the Yaml
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
				// Send all the newer posts to wallabag
				// fmt.Println("send to Bag:", element.Link)
				send(element.Link, tags)
			}
		}
		// Add a space for prettyness
		fmt.Println()

		// For the current feed update the latest post for the first url
		// received from the rss parser
		f.Feeds[i].LatestPost = feed.Items[0].GUID

		// Increase the count
		i++
	}

	// Update the configuration file
	viper.Set("feeds", f.Feeds)

	// Write the current viper config back to the config file
	// ( with updated latestposts )
	err = viper.WriteConfig()

	// Check for errors
	if err != nil {
		// Print the error to the console
		panic(err)
	}
}

func send(postURL string, Tags string) {
	// BagConfig : The Wallabag API will need a key and access token to work with
	type BagConfig struct {
        BaseUrl string `mapstructure:"baseUrl"`
		AccessToken string `mapstructure:"access_token"`
	}

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

    bearer := "Bearer " + Config.AccessToken

	// Configure a new request using the URL and Json
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
    req.Header.Add("Authorization", bearer)

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
