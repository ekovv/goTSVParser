package domains

import (
	"context"
	"goTSVParser/internal/shema"
)

//go:generate go run github.com/vektra/mockery/v3 --name=Service
type Service interface {
	Worker() error
	GetAll(ctx context.Context, r shema.Request) ([][]shema.Tsv, error)
}
