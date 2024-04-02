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

func (s *Service) Worker(ctx context.Context) error {
	const op = "service.Worker"

	checkedFiles, err := s.storage.GetCheckedFiles()
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : failed to get checked files: %v", op, err))
		return err
	}
	s.watcher.InitCheckedFiles(checkedFiles)

	out := make(chan string)
	go s.watcher.Scan(ctx, out)

loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case file, ok := <-out:
			if !ok {
				return nil
			}

			tsvChan, guidChan, errChan := s.parser.ParseFileAsync(file)
			var tsvArray []shema.Tsv
			var guidArray []string

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case tsv, ok := <-tsvChan:
					if !ok {
						tsvChan = nil
					} else {
						tsvArray = append(tsvArray, tsv)

						err = s.storage.Save(tsv)
						if err != nil {
							s.logger.Info(fmt.Sprintf("%s : failed to save data in db: %v", op, err))
							return err
						}
					}
				case guid, ok := <-guidChan:
					if !ok {
						guidChan = nil
					} else {
						guidArray = append(guidArray, guid)
					}
				case err, ok := <-errChan:
					if !ok {
						errChan = nil
					} else if err != nil {
						s.logger.Info(fmt.Sprintf("%s : failed to parse file: %v", op, err))

						f := shema.Files{
							File: file,
							Err:  err.Error(),
						}

						err = s.storage.SaveFilesWithErr(f)
						if err != nil {
							s.logger.Info(fmt.Sprintf("%s : failed to save file info in db: %v", op, err))
							return err
						}
						continue loop
					}
				}

				if tsvChan == nil && guidChan == nil && errChan == nil {
					break
				}
			}

			err = s.storage.SaveFiles(file)
			if err != nil {
				s.logger.Info(fmt.Sprintf("%s : failed to save file info in db: %v", op, err))
				return err
			}

			err = s.writer.WritePDF(tsvArray, guidArray)
			if err != nil {
				s.logger.Info(fmt.Sprintf("%s : failed to write pdf: %v", op, err))
				return err
			}
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
