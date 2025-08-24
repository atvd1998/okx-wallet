package service

type OKXRepository interface {
	GetAPIStatus() (bool, error)
}
