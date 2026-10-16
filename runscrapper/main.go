package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
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

func main() {
	db, err := clientDB()
	if err != nil {
		panic(err)
	}

	args := os.Args

	if len(args) == 2 && args[1] == "import" {
		runImport(db)
		return
	}

	runWeb(db)
}

func runWeb(db *sql.DB) {
	log.Println("Server starting on http://localhost:8080/")
	log.Fatal(serve(db))
}

func runImport(db *sql.DB) {
	client := &http.Client{}
	accessToken := getAccessToken(client)
	activities := getStravaActivities(client, accessToken)

	err := insertActivitiesIntoDB(context.Background(), db, activities)
	if err != nil {
		log.Fatalf("could not insert activities: %v", err)
	}
}
