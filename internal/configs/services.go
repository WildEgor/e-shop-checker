package configs

import (
	"encoding/json"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"log/slog"
	"os"
	"path/filepath"
)

const ConfigName = "services.json"

type ServiceUrl struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	Enabled bool   `json:"enabled"`
}

type ServicesConfig struct {
	Timeout int8         `json:"timeout"`
	URLs    []ServiceUrl `json:"urls"`
}

func NewServicesConfig() *ServicesConfig {
	var configPath string
	var config ServicesConfig

	dir, _ := os.Getwd()

	configPath = filepath.Join(dir, ConfigName)

	bts, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("cannot open file", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		panic(err)
	}

	if err := json.Unmarshal(bts, &config); err != nil {
		slog.Error("cannot parse file", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		panic(err)
	}

	return &config
}
