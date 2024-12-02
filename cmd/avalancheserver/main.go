package main

import (
	"avalancheserver/internal/aaa_api"
	apiv1 "avalancheserver/internal/api_v1"
	apiv2 "avalancheserver/internal/api_v2"
	apiv3 "avalancheserver/internal/api_v3"
	"avalancheserver/internal/config"
	"avalancheserver/internal/metrics"
	"flag"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func main() {
	configPath := flag.String("config", "", "path to JSON config file")
	flag.Parse()

	if configPath == nil || *configPath == "" {
		slog.Error("no_config_path_specified")
		os.Exit(1)
	}

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		slog.Error("error_loading_config_file", slog.Any("error", err))
		os.Exit(1)
	}

	apiRequester := aaa_api.APIRequester{}
	v1controller := apiv1.APIv1Controller{Requester: &apiRequester}
	v2controller := apiv2.APIv2Controller{Requester: &apiRequester}
	v3controller := apiv3.APIv3Controller{Requester: &apiRequester}

	r := gin.Default()

	if conf.GoatCounter.Enabled {
		r.Use(metrics.NewMiddleware(conf.GoatCounter.SiteCode, conf.GoatCounter.APIToken))
	}

	r.Use(sloggin.NewWithConfig(slog.Default(), sloggin.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      true,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  true,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,
	}))

	r.GET("/forecast/:center", v1controller.GetForecast)
	r.GET("/v2/forecast/:center", v2controller.GetForecast)
	r.GET("/v2/forecast/:center/:region", v2controller.GetForecast)
	r.GET("/v3/forecast/:center", v3controller.GetForecast)
	err = r.Run()
	if err != nil {
		panic(err)
	}
}
