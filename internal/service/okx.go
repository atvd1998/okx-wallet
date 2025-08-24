package service

import (
	"fmt"
	"okx-wallet/config"
	"okx-wallet/pkg/logger"
)

type OKXService struct {
	conf          *config.Config
	logger        *logger.Logger
	okxRepository OKXRepository
}

func NewOKXService(conf *config.Config, okxRepository OKXRepository) *OKXService {
	return &OKXService{
		conf:          conf,
		logger:        logger.MustNamed("okx_service"),
		okxRepository: okxRepository,
	}
}

func (s *OKXService) GetConnection() {
	test, err := s.okxRepository.GetAPIStatus()
	if err != nil {
		s.logger.Errorw("Failed to get API status", "error", err)
	}
	s.logger.Infow("API status", "status", fmt.Sprintf("%t", test))
}
