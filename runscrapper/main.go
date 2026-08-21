package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	StravaAPIURL   string
	StravaOAuthURL string
	ClientID       string
	ClientSecret   string
	RefreshToken   string
)

func init() {
	stravaAPIURL, ok := os.LookupEnv("STRAVA_API_URL")
	if ok {
		StravaAPIURL = stravaAPIURL
	}

	stravaOAuthURL, ok := os.LookupEnv("STRAVA_OAUTH_URL")
	if ok {
		StravaOAuthURL = stravaOAuthURL
	}

	clientID, ok := os.LookupEnv("CLIENT_ID")
	if ok {
		ClientID = clientID
	}

	clientSecret, ok := os.LookupEnv("CLIENT_SECRET")
	if ok {
		ClientSecret = clientSecret
	}

	refreshToken, ok := os.LookupEnv("REFRESH_TOKEN")
	if ok {
		RefreshToken = refreshToken
	}
}

func retrieveAccessToken(client *http.Client) string {
	data := url.Values{}
	data.Set("client_id", ClientID)
	data.Set("client_secret", ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", RefreshToken)

	req, err := http.NewRequest("POST", StravaOAuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Fatalf("API Error (Status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenData StravaTokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	fmt.Println("Success! New tokens received.")

	return tokenData.AccessToken
}

func main() {
	client := &http.Client{}
	accessToken := retrieveAccessToken(client)

	req, err := http.NewRequest("GET", StravaAPIURL+"/activities", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var activitiesData []StravaActivity
	if err = json.NewDecoder(resp.Body).Decode(&activitiesData); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	resp.Body.Close()

	for _, a := range activitiesData {
		fmt.Printf("%v\n", a)
	}

	// TODO: Export to DB
}
