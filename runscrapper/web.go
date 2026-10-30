package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates/*
var templateFS embed.FS

func importHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		client := &http.Client{}
		accessToken := getAccessToken(client)
		activities := getStravaActivities(client, accessToken)

		if err := insertActivitiesIntoDB(r.Context(), db, activities); err != nil {
			log.Printf("could not insert activities: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"ok": false, "error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(map[string]any{"ok": true, "activities": activities})
	}
}

func activitiesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT
			id_str, name, distance, moving_time, elapsed_time,
			total_elevation_gain, activity_type, sport_type, start_date_local,
			average_speed, max_speed, average_cadence, average_heart_rate,
			elev_high, elev_low, suffer_score
			FROM activities
			ORDER BY start_date_local DESC`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []StravaActivity
		for rows.Next() {
			var u StravaActivity
			if err := rows.Scan(
				&u.ID,
				&u.Name,
				&u.Distance,
				&u.MovingTime,
				&u.ElapsedTime,
				&u.TotalElevationGain,
				&u.Type,
				&u.SportType,
				&u.StartDateLocal,
				&u.AverageSpeed,
				&u.MaxSpeed,
				&u.AverageCadence,
				&u.AverageHeartrate,
				&u.ElevHigh,
				&u.ElevLow,
				&u.SufferScore,
			); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFS(templateFS, "templates/*.html")
		if err != nil {
			panic(err)
		}

		if err := tmpl.Execute(w, users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func serve(db *sql.DB) error {
	log.Println("Server starting on http://localhost:8080/")
	http.HandleFunc("/", activitiesHandler(db))
	http.HandleFunc("/import", importHandler(db))
	return http.ListenAndServe(":8080", nil)
}
