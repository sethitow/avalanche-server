package apiv2

import (
	"avalancheserver/internal/aaa_api"
	"avalancheserver/internal/aaa_api/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetForecastOffSeason(t *testing.T) {
	requester := mocks.NewRequester(t)
	controller := APIv2Controller{Requester: requester}

	file, err := os.Open("../aaa_api/testdata/offseason_sac.json")
	assert.Nil(t, err)
	defer file.Close()

	decoder := json.NewDecoder(file)
	aaaResponse := aaa_api.Root{}
	err = decoder.Decode(&aaaResponse)
	assert.Nil(t, err)
	requester.On("GetForecastByCenter", "SAC").Return(aaaResponse, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.GET("/forecast/:center", controller.GetForecast)
	c.Request, _ = http.NewRequest(http.MethodGet, "/forecast/SAC", bytes.NewBuffer([]byte("{}")))

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)

	responseDecoder := json.NewDecoder(w.Body)
	response := Envelope{}
	err = responseDecoder.Decode(&response)
	assert.Nil(t, err)

	assert.Equal(t,
		Response{DangerLevel: -1,
			TravelAdvice: "Watch for signs of unstable snow such as recent avalanches, cracking in the snow, and audible collapsing. Avoid traveling on or under similar slopes.",
			OffSeason:    true},
		response.Data)

}

func TestGetForecastAAAAPIError(t *testing.T) {
	requester := mocks.NewRequester(t)
	controller := APIv2Controller{Requester: requester}

	aaaResponse := aaa_api.Root{}
	requester.On("GetForecastByCenter", "SAC").Return(aaaResponse, aaa_api.ErrUpstreamAPIError)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.GET("/forecast/:center", controller.GetForecast)
	c.Request, _ = http.NewRequest(http.MethodGet, "/forecast/SAC", bytes.NewBuffer([]byte("{}")))

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	responseDecoder := json.NewDecoder(w.Body)
	response := EnvelopeError{}
	err := responseDecoder.Decode(&response)
	assert.Nil(t, err)
	assert.Equal(t, EnvelopeError{ResponseStatusError, "error requesting forecast"}, response)
}
