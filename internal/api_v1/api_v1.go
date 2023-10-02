package apiv1

import (
	"avalancheserver/internal/aaa_api"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RFC3339Z = "2006-01-02T15:04:05"
)

type Response struct {
	DangerLevel  int8   `json:"danger_level"`
	TravelAdvice string `json:"travel_advice"`
	UpdatedAt    int    `json:"updated_at"`
	ExpiresAt    int    `json:"expires_at"`
	OffSeason    bool   `json:"off_season"`
}

type APIv1Controller struct {
	Requester aaa_api.Requester
}

func (controller *APIv1Controller) GetForecast(c *gin.Context) {
	avalanche_center := c.Param("center")
	response, err := controller.Requester.GetForecastByCenter(avalanche_center)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(response.Features) < 1 {
		// Empty features means not center not found
		c.Status(http.StatusNotFound)
		return
	}

	if len(response.Features) > 1 {
		// V1 API only supports avalanche centers with one forecast area
		c.Status(http.StatusBadRequest)
		return
	}

	feature := response.Features[0]

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
			slog.String("start_date", feature.Properties.StartDate),
			slog.String("end_date", feature.Properties.EndDate))
		c.Status(http.StatusInternalServerError)
		return
	}

	var updatedAtInt = 0
	var expiresAtInt = 0

	if !errorParsingDates {
		updatedAtInt = int(updatedAt.Unix())
		expiresAtInt = int(expiresAt.Unix())
	}

	c.JSON(http.StatusOK, Response{
		DangerLevel:  int8(feature.Properties.DangerLevel),
		TravelAdvice: feature.Properties.TravelAdvice,
		UpdatedAt:    updatedAtInt,
		ExpiresAt:    expiresAtInt,
		OffSeason:    feature.Properties.OffSeason,
	})
}
