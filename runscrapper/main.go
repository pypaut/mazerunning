package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	StravaAPIURL   = ""
	StravaOAuthURL = ""
	ClientID       = ""
	ClientSecret   = ""
	RefreshToken   = ""
)

type StravaTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
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

	var activitiesData []Activity
	if err = json.NewDecoder(resp.Body).Decode(&activitiesData); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	resp.Body.Close()

	for _, a := range activitiesData {
		fmt.Printf("%v\n", a)
	}

	// TODO: Export to DB
}

type Activity struct {
	ID                 string  `json:"id_str"`
	Name               string  `json:"name"`
	Distance           float32 `json:"distance"`
	MovingTime         int     `json:"moving_time"`
	ElapsedTime        int     `json:"elapsed_time"`
	TotalElevationGain float32 `json:"total_elevation_gain"`
	Type               string  `json:"type"`
	SportType          string  `json:"sport_type"`
	StartDateLocal     string  `json:"start_date_local"`
	AverageSpeed       float32 `json:"average_speed"`
	MaxSpeed           float32 `json:"max_speed"`
	AverageCadence     float32 `json:"average_cadence"`
	AverageHeartrate   float32 `json:"average_heartrate"`
	ElevHigh           float32 `json:"elev_high"`
	ElevLow            float32 `json:"elev_low"`
	SufferScore        float32 `json:"suffer_score"`
}
