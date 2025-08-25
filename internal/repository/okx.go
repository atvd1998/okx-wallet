package repository

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"okx-wallet/config"
	"okx-wallet/pkg/logger"
	"time"

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
		SetBaseURL(conf.OKX.Web3Url).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetLogger(logger)
	return &OKXRepository{
		conf:   conf,
		logger: logger,
		client: client,
	}
}

func (r *OKXRepository) prepareRequest(method, endpoint string, body string) *resty.Request {
	// Generate timestamp for this specific request
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	req := r.client.R().
		SetHeader("OK-ACCESS-KEY", r.conf.OKX.APIKey).
		SetHeader("OK-ACCESS-SIGN", r.generateSignature(timestamp, method, endpoint, body)).
		SetHeader("OK-ACCESS-TIMESTAMP", timestamp).
		SetHeader("OK-ACCESS-PASSPHRASE", r.conf.OKX.PassPharse)

	return req
}

// generateSignature creates the required signature for OKX API authentication
func (r *OKXRepository) generateSignature(timestamp, method, requestPath, body string) string {
	// Create the prehash string
	prehash := timestamp + method + requestPath + body

	// Create the HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(r.conf.OKX.SecretKey))
	h.Write([]byte(prehash))

	// Return base64 encoded signature
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GetAPIStatus checks if the OKX API is available
func (r *OKXRepository) GetAPIStatus() (bool, error) {
	endpoint := "/api/v5/wallet/chain/supported-chains"
	res := make(map[string]any)
	httpResp, err := r.prepareRequest("GET", endpoint, "").
		SetResult(&res).
		Get(endpoint)

	r.logger.Infow("resp nè", "resp", httpResp)
	if err != nil {
		r.logger.Errorw("Failed to get API status", "error", err)
		return false, err
	}

	if httpResp.IsSuccess() {
		return true, nil
	}

	r.logger.Infow("resp nè", "resp", res)

	return false, fmt.Errorf("API returned non-success status: %s", httpResp.Status())
}
