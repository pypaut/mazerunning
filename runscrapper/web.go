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
    <title>Activities</title>
    <style>
        table { border-collapse: collapse; width: 50%; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <h2>Activities</h2>
    <table>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Distance</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.Name}}</td>
            <td>{{.Distance}}</td>
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
		rows, err := db.Query("SELECT id, name, distance FROM activities")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []StravaActivity
		for rows.Next() {
			var u StravaActivity
			if err := rows.Scan(&u.ID, &u.Name, &u.Distance); err != nil {
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
