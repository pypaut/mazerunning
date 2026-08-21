package main

type StravaActivity struct {
	ID                 string  `json:"id_str"`
	Name               string  `json:"name"`
	Distance           float32 `json:"distance"`
	MovingTime         int     `json:"moving_time"`
	ElapsedTime        int     `json:"elapsed_time"`
	TotalElevationGain float32 `json:"total_elevation_gain"`
	Type               string  `json:"type"`
	SportType          string  `json:"sport_type"`
	StartDateLocal     string  `json:"start_date_local"`
	AverageSpeed       float32 `json:"average_speed"`
	MaxSpeed           float32 `json:"max_speed"`
	AverageCadence     float32 `json:"average_cadence"`
	AverageHeartrate   float32 `json:"average_heartrate"`
	ElevHigh           float32 `json:"elev_high"`
	ElevLow            float32 `json:"elev_low"`
	SufferScore        float32 `json:"suffer_score"`
}

type StravaTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
