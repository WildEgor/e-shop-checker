package services

import (
	"context"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"log/slog"
	"net/http"
	"time"
)

const DefaultRequestTimeout = 5 * time.Second

type HttpCheckConfig struct {
	Sender  *Sender
	URL     string
	Timeout time.Duration
}

// NewHttpCheck New creates new HTTP service health check that verifies the following:
// - connection establishing
// - getting response status from defined URL
// - verifying that status code is above 400
func NewHttpCheck(cfg *HttpCheckConfig) func(ctx context.Context) error {
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultRequestTimeout
	}

	return func(ctx context.Context) error {
		req, err := http.NewRequest(http.MethodGet, cfg.URL, nil)
		if err != nil {
			slog.Error("creating the request for the health check failed", models.LogEntryAttr(&models.LogEntry{
				Err: err,
			}))
			return nil
		}

		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()

		// Inform remote service to close the connection after the transaction is complete
		req.Header.Set("Connection", "close")
		req = req.WithContext(ctx)

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			slog.Error("making the request for the health check failed", models.LogEntryAttr(&models.LogEntry{
				Err: err,
				Props: map[string]interface{}{
					"url": cfg.URL,
				},
			}))
			return nil
		}

		if resp != nil {
			defer resp.Body.Close()

			if resp.StatusCode >= http.StatusBadRequest {
				if cfg.Sender != nil {
					cfg.Sender.Send(SenderData{
						text: "Service <code>" + resp.Status + "</code> is down\n" + "Status: <b>" + cfg.URL + "</b>",
					})
				}

				slog.Error("service is not available at the moment", models.LogEntryAttr(&models.LogEntry{
					Props: map[string]interface{}{
						"url": cfg.URL,
					},
				}))
			} else {
				slog.Debug("service is available at the moment", models.LogEntryAttr(&models.LogEntry{
					Props: map[string]interface{}{
						"url": cfg.URL,
					},
				}))
			}
		}

		return nil
	}
}
