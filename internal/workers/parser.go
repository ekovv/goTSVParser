package workers

import (
	"encoding/csv"
	"goTSVParser/config"
	"goTSVParser/internal/constants"
	"goTSVParser/internal/shema"
	"io"
	"os"
	"strings"
)

type Parser struct {
	dirFrom string
}

func NewParser(cfg config.Config) *Parser {
	return &Parser{dirFrom: cfg.DirectoryFrom}
}

func (s *Parser) ParseFileAsync(fileName string) (<-chan shema.Tsv, <-chan string, <-chan error) {
	tsvChan := make(chan shema.Tsv)
	guidChan := make(chan string)
	errChan := make(chan error)

	go func() {
		defer close(tsvChan)
		defer close(guidChan)
		defer close(errChan)

		guidMap := make(map[string]bool)

		file, err := os.Open(fileName)
		if err != nil {
			errChan <- err
			return
		}

		if !strings.HasSuffix(file.Name(), ".tsv") {
			errChan <- constants.ErrNotTSV
			return
		}

		defer file.Close()

		reader := csv.NewReader(file)
		reader.Comma = '\t'
		for {
			str, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					return
				}
				errChan <- err
				return
			}
			if str == nil {
				break
			}
			if len(strings.TrimSpace(str[3])) < 10 {
				continue
			}

			t := shema.Tsv{
				Number:       strings.TrimSpace(str[0]),
				MQTT:         strings.TrimSpace(str[1]),
				InventoryID:  strings.TrimSpace(str[2]),
				UnitGUID:     strings.TrimSpace(str[3]),
				MessageID:    strings.TrimSpace(str[4]),
				MessageText:  strings.TrimSpace(str[5]),
				Context:      strings.TrimSpace(str[6]),
				MessageClass: strings.TrimSpace(str[7]),
				Level:        strings.TrimSpace(str[8]),
				Area:         strings.TrimSpace(str[9]),
				Address:      strings.TrimSpace(str[10]),
				Block:        strings.TrimSpace(str[11]),
				Type:         strings.TrimSpace(str[12]),
				Bit:          strings.TrimSpace(str[13]),
				InvertBit:    strings.TrimSpace(str[14]),
			}
			tsvChan <- t

			if _, exists := guidMap[t.UnitGUID]; !exists {
				guidChan <- t.UnitGUID
				guidMap[t.UnitGUID] = true
			}
		}
	}()

	return tsvChan, guidChan, errChan
}
