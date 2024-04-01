package domains

import (
	"context"
	"goTSVParser/internal/shema"
)

//go:generate go run github.com/vektra/mockery/v3 --name=Service
type Service interface {
	Scanner() error
	ParseFile(fileName string) ([]shema.Tsv, []string, error)
	WritePDF(tsv []shema.Tsv, unitGuid []string) error
	GetAll(ctx context.Context, r shema.Request) ([][]shema.Tsv, error)
}
