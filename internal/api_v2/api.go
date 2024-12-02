package apiv2

import (
	"avalancheserver/internal/aaa_api"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RFC3339Z = "2006-01-02T15:04:05"
)

type ResponseStatus string

const (
	ResponseStatusSuccess = ResponseStatus("success")
	ResponseStatusError   = ResponseStatus("error")
)

type Envelope struct {
	Status ResponseStatus `json:"status"`
	Data   Response       `json:"data"` // TODO: change to `any` if other data needs to be sent
}

type EnvelopeError struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
}

type Response struct {
	DangerLevel  int8   `json:"danger_level"`
	TravelAdvice string `json:"travel_advice"`
	UpdatedAt    int    `json:"updated_at"`
	ExpiresAt    int    `json:"expires_at"`
	OffSeason    bool   `json:"off_season"`
}

type APIv2Controller struct {
	Requester aaa_api.Requester
}

func (controller *APIv2Controller) GetForecast(c *gin.Context) {
	avalanche_center := c.Param("center")
	response, err := controller.Requester.GetForecastByCenter(avalanche_center)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, EnvelopeError{ResponseStatusError, "error requesting forecast"})
		return
	}

	if len(response.Features) < 1 {
		// Empty features means not center not found
		c.AbortWithStatusJSON(http.StatusNotFound, EnvelopeError{ResponseStatusError, "center not found"})
		return
	}

	var feature aaa_api.Feature
	regionStr := c.Param("region")
	if len(response.Features) == 1 {
		if regionStr != "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, EnvelopeError{ResponseStatusError, "region string specified, but center only has one region"})
			return
		}
		feature = response.Features[0]
	} else {
		if regionStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, EnvelopeError{ResponseStatusError, "center has multiple regions, but no region specified"})
			return
		}
		regionCode, err := strconv.ParseInt(regionStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, EnvelopeError{ResponseStatusError, "invalid region: could not parse"})
			return
		}
		for _, f := range response.Features {
			if f.ID == regionCode {
				feature = f
				break
			}
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, EnvelopeError{ResponseStatusError, "invalid region: region not found"})
		return
	}

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

	c.JSON(http.StatusOK, Envelope{
		ResponseStatusSuccess,
		Response{
			DangerLevel:  int8(feature.Properties.DangerLevel),
			TravelAdvice: feature.Properties.TravelAdvice,
			UpdatedAt:    updatedAtInt,
			ExpiresAt:    expiresAtInt,
			OffSeason:    feature.Properties.OffSeason,
		},
	})
}
