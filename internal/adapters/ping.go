package adapters

import (
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/configs"
	domains "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/domain"
)

type PingAdapter struct {
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

type IPingAdapter interface {
	GetApplicationStatus() *domains.StatusDomain
}

func NewPingAdapter(c *configs.AppConfig) *PingAdapter {
	return &PingAdapter{
		Version:     c.Version,
		Environment: c.GoEnv,
	}
}

func (pa *PingAdapter) GetApplicationStatus() *domains.StatusDomain {

	// TODO: check system health here

	return &domains.StatusDomain{
		Status:      "ok",
		Version:     pa.Version,
		Environment: pa.Environment,
	}
}
