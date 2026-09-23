package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.New("activities").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Activités</title>
    <style>
        table { border-collapse: collapse; width: 50%; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
	<a href="/import" style="appearance: button; text-decoration: none; padding: 10px; color: black; background-color: #f0f0f0; border: 1px solid #ccc; border-radius: 4px;">
		Importer les activités
	</a>
    <h2>Activities</h2>
    <table>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Distance</th>
			<th>MovingTime</th>
			<th>ElapsedTime</th>
			<th>TotalElevationGain</th>
			<th>Type</th>
			<th>SportType</th>
			<th>StartDateLocal</th>
			<th>AverageSpeed</th>
			<th>MaxSpeed</th>
			<th>AverageCadence</th>
			<th>AverageHeartrate</th>
			<th>ElevHigh</th>
			<th>ElevLow</th>
			<th>SufferScore</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.Name}}</td>
            <td>{{.Distance}}</td>
			<td>{{.MovingTime}}</td>
			<td>{{.ElapsedTime}}</td>
			<td>{{.TotalElevationGain}}</td>
			<td>{{.Type}}</td>
			<td>{{.SportType}}</td>
			<td>{{.StartDateLocal}}</td>
			<td>{{.AverageSpeed}}</td>
			<td>{{.MaxSpeed}}</td>
			<td>{{.AverageCadence}}</td>
			<td>{{.AverageHeartrate}}</td>
			<td>{{.ElevHigh}}</td>
			<td>{{.ElevLow}}</td>
			<td>{{.SufferScore}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>
`))

func importHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		accessToken := getAccessToken(client)
		activities := getStravaActivities(client, accessToken)

		err := insertActivitiesIntoDB(r.Context(), db, activities)
		if err != nil {
			log.Fatalf("could not insert activities: %v", err)
		}
	}
}

func activitiesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT
			id_str, name, distance, moving_time, elapsed_time,
			total_elevation_gain, activity_type, sport_type, start_date_local,
			average_speed, max_speed, average_cadence, average_heart_rate,
			elev_high, elev_low, suffer_score
			FROM activities`)
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

		if err := tmpl.Execute(w, users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func serve(db *sql.DB) error {
	http.HandleFunc("/activities", activitiesHandler(db))
	http.HandleFunc("/import", importHandler(db))
	return http.ListenAndServe(":8080", nil)
}
