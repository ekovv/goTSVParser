package domains

import "goTSVParser/internal/shema"

type Storage interface {
	SaveFiles(sh shema.Files) error
}
