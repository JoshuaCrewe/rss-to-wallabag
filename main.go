package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/viper"
    "io"
    "log"
	"net/http"
    // "net/http/httputil"
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

type BagConfig struct {
    BaseUrl string `mapstructure:"baseUrl"`
    AccessToken string `mapstructure:"access_token"`
    ClientID string `mapstructure:"client_id"`
    ClientSecret string `mapstructure:"client_secret"`
    UserName string `mapstructure:"username"`
    Password string `mapstructure:"password"`
}

type Response struct {
    AccessToken string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// @TODO Delete this when done
func debug(data []byte, err error) {
    if err == nil {
        fmt.Printf("%s\n\n", data)
    } else {
        log.Fatalf("%s\n\n", err)
    }
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

func main() {
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

func send(postURL string, Tags string, AccessToken string) {

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
