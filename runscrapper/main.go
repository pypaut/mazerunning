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
	activities := getStravaActivities(client, accessToken, knownActivityIDs(db))

	err := insertActivitiesIntoDB(context.Background(), db, activities)
	if err != nil {
		log.Fatalf("could not insert activities: %v", err)
	}
}

// knownActivityIDs returns the set of id_str values already stored, so the
// fetch loop can stop paging once it reaches already-synced territory instead
// of re-walking (and re-inserting) the full history on every sync.
func knownActivityIDs(db *sql.DB) map[string]bool {
	rows, err := db.Query("SELECT id_str FROM activities")
	if err != nil {
		log.Printf("could not load known ids: %v — treating DB as empty", err)
		return map[string]bool{}
	}
	defer rows.Close()
	ids := map[string]bool{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids[id] = true
		}
	}
	return ids
}
