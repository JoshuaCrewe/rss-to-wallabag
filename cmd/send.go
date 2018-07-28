// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
)

// runCmd represents the run command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Process RSS feeds and send them to pocket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processFeeds()
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type feedItem struct {
	URL        string `mapstructure:"url"`
	Tags       string `mapstructure:"tags"`
	LatestPost string `mapstructure:"latestpost"`
}

type feeds struct {
	Feeds []feedItem
}

func processFeeds() {
	// Create a new struct for feeds
	f := feeds{}

	// Populate this structure with data from the yaml
	err := viper.Unmarshal(&f)

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
				// Send all the newer posts to pocket
				// fmt.Println("send to Pocket:", element.Link)

				sendToPocket(element.Link, tags)
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
		fmt.Println("there was an error with writing the config ", err)
	}
}

func sendToPocket(postURL string, Tags string) {
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

	// Populate struct with config data
	err := viper.Unmarshal(&Config)

	// Check for errors
	if err != nil {
		// Print the error to the console
		fmt.Println("there was an error Unmarshalling ", err)
	}

	// Pocket API enpoint for adding items
	baseURL := "https://getpocket.com/v3/add"

	// Encode string as a URL
	// https://getpocket.com/developer/docs/v3/add - Best Practices
	u, err := url.Parse(postURL)

	// Handle errors
	if err != nil {
		fmt.Println("there was an error with the http request", err)
		return
	}
	// Convert encoded URL back to a string
	URL := u.String()

	// Gather data use use in POST request
	jsonStr := &PostRequest{URL, Config.ConsumerKey, Config.AccessToken, Tags}

	// Json encode this data
	b, err := json.Marshal(jsonStr)
	// Handle errors
	if err != nil {
		fmt.Println("there was an error with the http request", err)
		return
	}

	// Configure a new request using the URL and Json
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(b))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	// Handle errors
	if err != nil {
		fmt.Println("there was an error with the http request", err)
		return
	}

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)

	// Handle errors
	if err != nil {
		fmt.Println("there was an error with the http request", err)
		return
	}

	// Close the response
	defer resp.Body.Close()

	// Inform the user which URL has been sent
	fmt.Println("🚀  ", URL)
}
