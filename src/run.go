package run

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/viper"
)

type Response struct {
    AccessToken string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

type feedItem struct {
	URL        string `mapstructure:"url"`
	Tags       string `mapstructure:"tags"`
	LatestPost string `mapstructure:"latestpost"`
}

type feeds struct {
	Feeds []feedItem
}

func run() {
    viper.SetConfigType("yaml")

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

    auth := auth()
    var response Response

    err = json.Unmarshal([]byte(auth), &response)

    fmt.Println("🍺  response ", response.AccessToken)

	// Get the pocket specific data
	Config := BagConfig{}

	// Populate struct with config data
	err = viper.Unmarshal(&Config)

	// Check for errors
	if err != nil {
		// Print the error to the console
		panic(err)
	}

    // Config.AccessToken = response.AccessToken

    // Update the configuration file
	// viper.Set("response.AccessToken", Config.AccessToken)

    // fmt.Println("Config after auth", Config)
	// viper.AddConfigPath("$HOME/.config/bag")
    // viper.SafeWriteConfig()
    // return

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
	// fmt.Println("🐛  ", f.Feeds)

	// For each of the feeds from the Yaml
	for _, element := range f.Feeds {

		// Indicate something is happeing
		fmt.Println("🍺  downloading ", element.URL)

		// parse the RSS url for that feed
		feed, err := fp.ParseURL(element.URL)

        if err != nil {
            // Print the error to the console
            // fmt.Println("🐛  ", err)
            // fmt.Println(" ")
            i++
            // return
            continue
        }

        // if err == nil {
            

		// Save the ID for the last post fetched
		latest := element.LatestPost

		// Save the tags from the yaml
		tags := element.Tags

		// For each of the items in that feed
		for _, element := range feed.Items {

			// If the current ID is the last one we got
			if element.Link == latest {
				// Stop looking through the posts
				break
			} else {
				// Send all the newer posts to wallabag
				// fmt.Println("Element:", (element))
				// fmt.Println("send to Bag:", element.Link)
				// fmt.Println("Current:", element.GUID)
				// fmt.Println("send to Bag:", tags)
				send(element.Link, tags, response.AccessToken)
			}
		}
		// Add a space for prettyness
		fmt.Println()

        if len(feed.Items) > 0 {
            // For the current feed update the latest post for the first url
            // received from the rss parser
            f.Feeds[i].LatestPost = feed.Items[0].Link
        }

        // Update the configuration file
        viper.Set("feeds", f.Feeds)

        // Write the current viper config back to the config file
        // ( with updated latestposts )
        err = viper.WriteConfig()

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
