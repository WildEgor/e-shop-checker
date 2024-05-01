package services

import (
	"context"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/adapters"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/configs"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// CheckerService Service Facade contains logic for check services and send notifications
type CheckerService struct {
	notificator        adapters.INotificator
	healthCheckAdapter *adapters.HealthCheckAdapter
	servicesConfig     *configs.ServicesConfig
}

func NewCheckerService(
	notificator adapters.INotificator,
	healthCheckAdapter *adapters.HealthCheckAdapter,
	servicesConfig *configs.ServicesConfig,
) *CheckerService {
	return &CheckerService{
		notificator:        notificator,
		healthCheckAdapter: healthCheckAdapter,
		servicesConfig:     servicesConfig,
	}
}

// Check simple running in goroutine
func (s *CheckerService) Check() {
	log.Debug("check urls", models.LogEntryAttr(&models.LogEntry{
		Props: map[string]interface{}{
			"urls": s.servicesConfig.URLs,
		},
	}))

	for {
		sleep := time.Duration(s.servicesConfig.Timeout)

		for i := 0; i < len(s.servicesConfig.URLs); i++ {
			if !s.servicesConfig.URLs[i].Enabled {
				continue
			}

			func() {
				log.Debug("Start check...")

				resp, err := http.Get(s.servicesConfig.URLs[i].URL)
				defer resp.Body.Close()

				if err != nil {
					log.Warn("Error connect to server: " + s.servicesConfig.URLs[i].URL)
					return
				}

				log.Debug("Connect to server: " + s.servicesConfig.URLs[i].URL + " OK!")

				if resp.StatusCode != 200 {
					if err := s.notificator.Send("Service <code>" + resp.Status + "</code> is down\n" + "Status: <b>" + s.servicesConfig.URLs[i].URL + "</b>"); err != nil {
						log.Warn("Cannot send telegram alert. Reason: ", err)
					}

					sleep = time.Duration(s.servicesConfig.Timeout)
				}
			}()
		}

		time.Sleep(sleep * time.Second)
	}
}

// ServicesCheck more advanced use case
func (s *CheckerService) ServicesCheck(ctx context.Context) {
	for i := 0; i < len(s.servicesConfig.URLs); i++ {
		if !s.servicesConfig.URLs[i].Enabled {
			continue
		}

		log.Info("Checking service: ", s.servicesConfig.URLs[i].URL)

		s.healthCheckAdapter.Register(adapters.HealthConfig{
			Name:      s.servicesConfig.URLs[i].ID,
			Timeout:   time.Duration(s.servicesConfig.Timeout),
			SkipOnErr: false,
			Check: NewHttpCheck(&HttpCheckConfig{
				Sender: InitSender(s.notificator),
				URL:    s.servicesConfig.URLs[i].URL,
			}),
		})
	}

	for {
		s.healthCheckAdapter.Measure(ctx)

		time.Sleep(time.Duration(5) * time.Second)
	}
}
