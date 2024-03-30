package service

import (
	"go.uber.org/zap"
	"goTSVParser/config"
	"goTSVParser/internal/domains"
	"goTSVParser/internal/watcher"
)

type Service struct {
	storage domains.Storage
	watcher *watcher.Watcher
	config  config.Config
	logger  *zap.Logger
}

func NewService(storage domains.Storage, watcher *watcher.Watcher, config config.Config) *Service {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil
	}
	return &Service{storage: storage, watcher: watcher, config: config, logger: logger}
}

func (s *Service) Scanner() error {
	const op = "service.Scanner"
	out := make(chan string)
	go s.watcher.Scan(out)
	for file := range out {

	}
	return nil
}
