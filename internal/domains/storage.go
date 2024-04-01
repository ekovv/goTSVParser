package domains

import (
	"context"
	"goTSVParser/internal/shema"
)

//go:generate go run github.com/vektra/mockery/v3 --name=Storage
type Storage interface {
	SaveFiles(sh shema.Files) error
	Save(sh shema.Tsv) error
	GetCheckedFiles() ([]shema.ParsedFiles, error)
	GetAllGuids(ctx context.Context, unitGuid string) ([]shema.Tsv, error)
	ShutDown() error
}
