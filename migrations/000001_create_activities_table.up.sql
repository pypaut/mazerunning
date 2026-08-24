CREATE TABLE if NOT EXISTS activities (
    id SERIAL PRIMARY KEY,
    id_str VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    distance NUMERIC NOT NULL,
    moving_time INTEGER NOT NULL,
    elapsed_time INTEGER NOT NULL,
    total_elevation_gain NUMERIC NOT NULL,
    activity_type VARCHAR(255) NOT NULL,
    sport_type VARCHAR(255) NOT NULL,
    start_date_local VARCHAR(255) NOT NULL,
    average_speed NUMERIC NOT NULL,
    max_speed NUMERIC NOT NULL,
    average_cadence NUMERIC NOT NULL,
    average_heart_rate NUMERIC NOT NULL,
    elev_high NUMERIC NOT NULL,
    elev_low NUMERIC NOT NULL,
    suffer_score NUMERIC NOT NULL
)
