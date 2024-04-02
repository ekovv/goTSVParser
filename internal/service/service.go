package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"goTSVParser/config"
	"goTSVParser/internal/constants"
	"goTSVParser/internal/domains"
	"goTSVParser/internal/shema"
	"goTSVParser/internal/workers"
)

type Service struct {
	storage domains.Storage
	watcher *workers.Watcher
	parser  *workers.Parser
	writer  *workers.Writer
	config  config.Config
	logger  *zap.Logger
}

func NewService(storage domains.Storage, watcher *workers.Watcher, parser *workers.Parser, writer *workers.Writer, config config.Config) *Service {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil
	}
	return &Service{storage: storage, watcher: watcher, config: config, logger: logger, parser: parser, writer: writer}
}

func (s *Service) Worker() error {
	const op = "service.Worker"

	checkedFiles, err := s.storage.GetCheckedFiles()
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : %v", op, err))
		return fmt.Errorf("failed to get checked files")
	}
	s.watcher.InitCheckedFiles(checkedFiles)

	out := make(chan string)
	go s.watcher.Scan(out)

	for file := range out {
		tsv, unitGuid, err := s.parser.ParseFile(file)
		var errFrom string
		if err != nil {
			s.logger.Info(fmt.Sprintf("%s : %v", op, err))
			errFrom = err.Error()
		}

		f := shema.Files{
			File: file,
			Err:  errFrom,
		}
		err = s.storage.SaveFiles(f)
		if err != nil {
			s.logger.Info(fmt.Sprintf("%s : %v", op, err))
			return err
		}

		for _, ts := range tsv {
			err = s.storage.Save(ts)
			if err != nil {
				s.logger.Info(fmt.Sprintf("%s : %v", op, err))
				return fmt.Errorf("failed to save data in db: %w", err)
			}
		}

		err = s.writer.WritePDF(tsv, unitGuid)
		if err != nil {
			s.logger.Info(fmt.Sprintf("%s : %v", op, err))
			return fmt.Errorf("failed to write pdf: %w", err)
		}
	}

	return nil
}

func (s *Service) GetAll(ctx context.Context, r shema.Request) ([][]shema.Tsv, error) {
	const op = "service.GetAll"

	tsvFromDB, err := s.storage.GetAllGuids(ctx, r.UnitGUID)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : %v", op, err))
		return nil, constants.ErrNotFound
	}
	var resultArray [][]shema.Tsv

	arrayWithPage := SubArray(r.Page, tsvFromDB)
	for i := 0; i < len(arrayWithPage); i += r.Limit {
		end := i + r.Limit

		if end > len(arrayWithPage) {
			end = len(arrayWithPage)
		}

		resultArray = append(resultArray, arrayWithPage[i:end])
	}

	return resultArray, nil
}

func SubArray(startIndex int, data []shema.Tsv) []shema.Tsv {
	if startIndex < 0 || startIndex >= len(data) {
		return nil
	}

	return data[startIndex:]
}
