package controller

import (
	"okx-wallet/config"
	"okx-wallet/pkg/logger"
)

type OKXService interface {
	GetConnection()
}

type Controller struct {
	conf       *config.Config
	okxService OKXService
	logger     *logger.Logger
}

func NewController(conf *config.Config, okxService OKXService) *Controller {
	return &Controller{
		okxService: okxService,
		conf:       conf,
		logger:     logger.MustNamed("controller"),
	}
}
