package service

import "fmt"

type OKXService struct{}

func NewOKXService() *OKXService {
	return &OKXService{}
}

func (s *OKXService) GetConnection() {
	fmt.Println("test service")
}
