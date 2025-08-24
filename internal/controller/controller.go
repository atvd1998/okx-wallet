package controller

import "okx-wallet/config"

type OKXService interface {
	GetConnection()
}

type Controller struct {
	conf       *config.Config
	okxService OKXService
}

func NewController(conf *config.Config, okxService OKXService) *Controller {
	return &Controller{
		okxService: okxService,
		conf:       conf,
	}
}
