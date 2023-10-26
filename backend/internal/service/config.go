package service

type Config interface {
	Load() error
}
