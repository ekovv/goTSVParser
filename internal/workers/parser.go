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

func (s *Parser) ParseFile(fileName string) ([]shema.Tsv, []string, error) {
	file, err := os.Open(s.dirFrom + "/" + fileName)
	if err != nil {
		return nil, nil, err
	}

	if !strings.HasSuffix(file.Name(), ".tsv") {
		return nil, nil, constants.ErrNotTSV
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	var data []shema.Tsv
	var array []string
	for {
		for _, d := range data {
			if array == nil {
				array = append(array, d.UnitGUID)
			}
		}
		str, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return data, array, nil
			}
			return nil, nil, err
		}
		if str == nil {
			break
		}
		if len(strings.TrimSpace(str[3])) < 10 {
			continue
		}
	loop:
		for _, s := range data {
			for _, guid := range array {
				if s.UnitGUID == guid || guid == strings.TrimSpace(str[3]) {
					continue loop
				}
			}
			array = append(array, strings.TrimSpace(str[3]))

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
		data = append(data, t)
	}
	return data, array, nil
}
