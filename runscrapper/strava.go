package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	StravaAPIURL   string
	StravaOAuthURL string
	ClientID       string
	ClientSecret   string
	RefreshToken   string
)

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

func getStravaActivities(
	client *http.Client, accessToken string, knownIDs map[string]bool,
) (all []StravaActivity) {
	// Page through Strava's /activities newest→oldest at the largest page size.
	// Stop early when an entire page is already in the DB: since results are
	// newest-first, hitting a fully-known page means everything older is synced
	// too, so there's nothing left to backfill. This keeps steady-state syncs at
	// one page while still recovering missing history on a fresh/partial DB.
	const perPage = 200
	for page := 1; ; page++ {
		u := fmt.Sprintf("%s/activities?per_page=%d&page=%d", StravaAPIURL, perPage, page)

		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		var pageData []StravaActivity
		if err = json.NewDecoder(resp.Body).Decode(&pageData); err != nil {
			resp.Body.Close()
			log.Fatalf("Failed to decode JSON: %v", err)
		}
		resp.Body.Close()

		fmt.Printf("page %d: %d activities (status %s)\n", page, len(pageData), resp.Status)

		// Rate-limit guard: Strava returns "15min,daily" usage/limit. Stop before
		// burning the 15-min quota so a long backfill doesn't lock us out.
		if usage, limit, ok := parseRateLimit(resp.Header); ok && usage >= limit*9/10 {
			log.Printf("strava rate limit at %d/%d — pausing sync", usage, limit)
			break
		}

		// Whole page already in DB → everything older is synced too. Stop.
		if len(pageData) > 0 && allKnown(pageData, knownIDs) {
			break
		}

		all = append(all, pageData...)

		// Break on an empty page (end of history) or a hard page cap.
		if len(pageData) == 0 || page >= 50 {
			break
		}
	}

	return
}

func allKnown(page []StravaActivity, known map[string]bool) bool {
	for _, a := range page {
		if !known[a.ID] {
			return false
		}
	}
	return true
}

// parseRateLimit reads Strava's X-RateLimit-{Usage,Limit} headers, which are
// formatted as "15min,daily" (e.g. "300,1000"). Returns the 15-min figures.
func parseRateLimit(h http.Header) (usage, limit int, ok bool) {
	uParts := strings.Split(h.Get("X-RateLimit-Usage"), ",")
	lParts := strings.Split(h.Get("X-RateLimit-Limit"), ",")
	if len(uParts) == 0 || len(lParts) == 0 {
		return 0, 0, false
	}
	usage, errU := strconv.Atoi(uParts[0])
	limit, errL := strconv.Atoi(lParts[0])
	if errU != nil || errL != nil {
		return 0, 0, false
	}
	return usage, limit, true
}
