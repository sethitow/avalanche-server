package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagikazarmark/slog-shim"
)

// NewMiddleware creates a middleware function to send metrics asynchronously.
//
// token=[your api token]
// api=https://avalanche-sethitow-com.goatcounter.com/api/v0
//
//	curl -X POST "$api/count" \
//	    -H 'Content-Type: application/json' \
//	    -H "Authorization: Bearer $token" \
//	    --data '{"no_sessions": true, "hits": [{"path": "/one"}, {"path": "/two"}]}'
func NewMiddleware(siteCode string, apiToken string) gin.HandlerFunc {
	metricsURL := fmt.Sprintf("https://%s.goatcounter.com/api/v0/count", siteCode)
	return func(c *gin.Context) {
		countReq := APICountRequest{
			Hits: []APICountRequestHit{
				{
					IP:        c.ClientIP(),
					UserAgent: c.Request.UserAgent(),
					Path:      c.Request.URL.Path,
					Query:     c.Request.URL.RawQuery,
				},
			},
		}

		// Send metrics asynchronously.
		go func(data APICountRequest) {
			jsonData, err := json.Marshal(data)
			if err != nil {
				slog.Warn("failed_to_marshal_metrics_data",
					slog.Any("data", data),
					slog.Any("error", err))
				return
			}

			req, err := http.NewRequest("POST", metricsURL, bytes.NewBuffer(jsonData))
			if err != nil {
				slog.Warn("failed_to_instantiate_request", slog.Any("error", err))
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+apiToken)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				slog.Warn("failed_to_send_metrics", slog.Any("error", err))
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusAccepted {
				slog.Warn("unexpected_response_status",
					slog.String("status", resp.Status),
					slog.Int("status_code", resp.StatusCode))
			}
		}(countReq)

		c.Next()
	}
}
