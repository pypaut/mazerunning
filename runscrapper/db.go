package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Registers "pgx" driver
)

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

func insertActivitiesIntoDB(
	ctx context.Context, db *sql.DB, activities []StravaActivity,
) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begintx: %w", err)
	}
	defer tx.Rollback()

	// Insert into db
	query := `INSERT INTO activities (` +
		`id_str,` +
		`name,` +
		`distance,` +
		`moving_time,` +
		`elapsed_time,` +
		`total_elevation_gain,` +
		`activity_type,` +
		`sport_type,` +
		`start_date_local,` +
		`average_speed,` +
		`max_speed,` +
		`average_cadence,` +
		`average_heart_rate,` +
		`elev_high,` +
		`elev_low,` +
		`suffer_score) ` +
		`VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("could not prepare context: %w", err)
	}
	defer stmt.Close()

	for _, a := range activities {
		if _, err := stmt.ExecContext(
			ctx,
			a.ID,
			a.Name,
			a.Distance,
			a.MovingTime,
			a.ElapsedTime,
			a.TotalElevationGain,
			a.Type,
			a.SportType,
			a.StartDateLocal,
			a.AverageSpeed,
			a.MaxSpeed,
			a.AverageCadence,
			a.AverageHeartrate,
			a.ElevHigh,
			a.ElevLow,
			a.SufferScore,
		); err != nil {
			return fmt.Errorf("could not exec context: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction")
	}

	return nil
}
