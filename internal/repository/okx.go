package repository

import (
	"fmt"
	"okx-wallet/config"
	"okx-wallet/pkg/logger"

	"github.com/go-resty/resty/v2"
)

type OKXRepository struct {
	conf   *config.Config
	logger *logger.Logger
	client *resty.Client
}

func NewOKXRepository(conf *config.Config) *OKXRepository {
	logger := logger.MustNamed("okx_repository")
	client := resty.New()
	client.
		SetBaseURL(conf.OKX.Url).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetLogger(logger)
	return &OKXRepository{
		conf:   conf,
		logger: logger,
		client: client,
	}
}

// GetAPIStatus checks if the OKX API is available
func (r *OKXRepository) GetAPIStatus() (bool, error) {
	resp, err := r.client.R().
		Get("/api/v5/public/status")

	r.logger.Infow("resp nè", "resp", resp)
	if err != nil {
		r.logger.Errorw("Failed to get API status", "error", err)
		return false, err
	}

	if resp.IsSuccess() {
		return true, nil
	}

	return false, fmt.Errorf("API returned non-success status: %s", resp.Status())
}
