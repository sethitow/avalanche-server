package apiv3

import (
	"avalancheserver/internal/aaa_api"
	"net/http"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type ResponseStatus string

const (
	ResponseStatusSuccess = ResponseStatus("success")
	ResponseStatusError   = ResponseStatus("error")
)

const (
	RFC3339Z = "2006-01-02T15:04:05"
)

type Envelope[T any] struct {
	Status ResponseStatus `json:"status"`
	Data   T              `json:"data"`
}

type EnvelopeError struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
}

type Response struct {
	MostSevereDangerLevel int            `json:"most_severe_danger_level"`
	MostSevereAreaName    string         `json:"most_severe_area_name"`
	Areas                 []ForecastArea `json:"areas"`
}

type ForecastArea struct {
	Name         string `json:"name"`
	DangerLevel  int    `json:"danger_level"`
	TravelAdvice string `json:"travel_advice"`
	UpdatedAt    int    `json:"updated_at"`
	ExpiresAt    int    `json:"expires_at"`
	OffSeason    bool   `json:"off_season"`
}

type APIv3Controller struct {
	Requester aaa_api.Requester
}

func (controller *APIv3Controller) GetForecast(c *gin.Context) {
	avalanche_center := c.Param("center")
	response, err := controller.Requester.GetForecastByCenter(avalanche_center)
	if err != nil {
		slog.Error("error_getting_forecast_from_aaa_api", slog.Any("err", err))
		c.JSON(http.StatusInternalServerError, EnvelopeError{
			Status:  ResponseStatusError,
			Message: "error from Avalanche.org",
		})
		return
	}

	if len(response.Features) < 1 {
		// Empty features means not center not found
		c.JSON(http.StatusNotFound, EnvelopeError{
			Status:  ResponseStatusError,
			Message: "avalanche center not found",
		})
		return
	}

	forecastAreas := make([]ForecastArea, 0, len(response.Features))

	mostSevereDangerLevel := -1
	mostSevereAreaName := ""

	for _, feature := range response.Features {
		errorParsingDates := false
		updatedAt, err := time.Parse(RFC3339Z, feature.Properties.StartDate)
		if err != nil {
			slog.Error("error_parsing_end_date", slog.Any("err", err))
			errorParsingDates = true
		}
		expiresAt, err := time.Parse(RFC3339Z, feature.Properties.EndDate)
		if err != nil {
			slog.Error("error_parsing_end_date", slog.Any("err", err))
			errorParsingDates = true
		}

		if errorParsingDates && !feature.Properties.OffSeason {
			slog.Error("invalid_dates",
				slog.String("area_name", feature.Properties.Name),
				slog.String("start_date", feature.Properties.StartDate),
				slog.String("end_date", feature.Properties.EndDate))
			c.JSON(http.StatusInternalServerError, EnvelopeError{
				Status:  ResponseStatusError,
				Message: "forecast has invalid expiration date",
			})
			return
		}

		var updatedAtInt = 0
		var expiresAtInt = 0

		if !errorParsingDates {
			updatedAtInt = int(updatedAt.Unix())
			expiresAtInt = int(expiresAt.Unix())
		}
		forecastArea := ForecastArea{
			Name:         feature.Properties.Name,
			DangerLevel:  feature.Properties.DangerLevel,
			TravelAdvice: feature.Properties.TravelAdvice,
			UpdatedAt:    updatedAtInt,
			ExpiresAt:    expiresAtInt,
			OffSeason:    feature.Properties.OffSeason,
		}

		if forecastArea.DangerLevel > mostSevereDangerLevel {
			mostSevereAreaName = forecastArea.Name
			mostSevereDangerLevel = forecastArea.DangerLevel
		} else if mostSevereDangerLevel > -1 && forecastArea.DangerLevel == mostSevereDangerLevel {
			mostSevereAreaName = "Multiple"
		}

		forecastAreas = append(forecastAreas, forecastArea)
	}

	c.JSON(http.StatusOK, Envelope[Response]{
		Status: ResponseStatusSuccess,
		Data: Response{
			MostSevereDangerLevel: mostSevereDangerLevel,
			MostSevereAreaName:    mostSevereAreaName,
			Areas:                 forecastAreas,
		},
	},
	)
}
