package main

import (
	"avalancheserver/internal/aaa_api"
	apiv1 "avalancheserver/internal/api_v1"
	apiv2 "avalancheserver/internal/api_v2"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	apiRequester := aaa_api.APIRequester{}
	v1controller := apiv1.APIv1Controller{Requester: &apiRequester}
	v2controller := apiv2.APIv2Controller{Requester: &apiRequester}
	r.GET("/forecast/:center", v1controller.GetForecast)
	r.GET("/v2/forecast/:center", v2controller.GetForecast)
	r.Run()
}
