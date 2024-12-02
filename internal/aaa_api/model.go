package aaa_api

type Root struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string     `json:"type"`
	ID         int64      `json:"id"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	Name         string `json:"name"`
	DangerLevel  int    `json:"danger_level"`
	TravelAdvice string `json:"travel_advice"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	OffSeason    bool   `json:"off_season"`
}
