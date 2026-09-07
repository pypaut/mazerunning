package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Registers "pgx" driver
)

// User represents a table row
type User struct {
	ID       int
	Name     string
	Distance float32
}

// HTML template with standard Go template syntax
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

func usersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, distance FROM activities")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
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

func main() {
	db, err := clientDB()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/activities", usersHandler(db))

	log.Println("Server starting on http://localhost:8080/activities")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
