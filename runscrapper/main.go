package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // Registers "postgres" driver
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

func getAccessToken(client *http.Client) string {
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

func getStravaActivities(client *http.Client, accessToken string) (activitiesData []StravaActivity) {
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

	fmt.Printf("Response status: %s\n", resp.Status)

	if err = json.NewDecoder(resp.Body).Decode(&activitiesData); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	resp.Body.Close()

	return
}

func clientDB() (*sql.DB, error) {
	dbHost, exists := os.LookupEnv("DB_HOST")
	if !exists {
		log.Fatalf("could not get db host")
	}

	dbPort, exists := os.LookupEnv("DB_PORT")
	if !exists {
		log.Fatalf("could not get db port")
	}

	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		log.Fatalf("could not get db name")
	}

	dbUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		log.Fatalf("could not get db user")
	}

	dbPassword, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		log.Fatalf("could not get postgres password")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	return sql.Open("pgx", connStr)
}

func main() {
	// Get some activities
	client := &http.Client{}
	accessToken := getAccessToken(client)
	activities := getStravaActivities(client, accessToken)

	for _, a := range activities {
		fmt.Printf("%v\n", a)
	}

	// Connext to db
	db, err := clientDB()
	if err != nil {
		log.Fatalf("Erreur d'ouverture de la connexion : %v", err)
	}
	defer db.Close()
}
