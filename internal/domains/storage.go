package domains

import (
	"context"
	"goTSVParser/internal/shema"
)

type Storage interface {
	SaveFiles(sh shema.Files) error
	Save(sh shema.Tsv) error
	GetCheckedFiles() ([]shema.ParsedFiles, error)
	GetAllGuids(ctx context.Context, unitGuid string) ([]shema.Tsv, error)
	ShutDown() error
}
