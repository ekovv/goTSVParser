package workers

import (
	"encoding/csv"
	"errors"
	"goTSVParser/config"
	"goTSVParser/internal/constants"
	"goTSVParser/internal/shema"
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestService_ParseFile(t *testing.T) {
	type args struct {
		file string
		dir  string
	}
	tests := []struct {
		name      string
		args      args
		wantTsv   []shema.Tsv
		wantGuids []string
		wantErr   error
	}{
		{
			name: "OK#1",
			args: args{dir: "testDirectory", file: "OK1.tsv"},
			wantTsv: []shema.Tsv{
				{
					Number:       "5",
					InventoryID:  "G-044325",
					UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
					MessageID:    "cold78_Defrost_status",
					MessageText:  "Разморозка",
					MessageClass: "waiting",
					Level:        "100",
					Area:         "LOCAL",
					Address:      "cold78_status.Defrost_status",
				},
			},
			wantGuids: []string{
				"01749246-9617-585e-9e19-157ccad61ee2",
			},
			wantErr: nil,
		},
		{
			name: "OK#2",
			args: args{dir: "testDirectory", file: "OK2.tsv"},
			wantTsv: []shema.Tsv{
				{
					Number:       "5",
					InventoryID:  "G-044325",
					UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
					MessageID:    "cold78_Defrost_status",
					MessageText:  "Разморозка",
					MessageClass: "waiting",
					Level:        "100",
					Area:         "LOCAL",
					Address:      "cold78_status.Defrost_status",
				},
				{
					Number:       "6",
					InventoryID:  "G-044325",
					UnitGUID:     "01749246-9617-585e-9e19-157ccad61ee2",
					MessageID:    "cold78_VentSK_status",
					MessageText:  "Вентилятор",
					MessageClass: "working",
					Level:        "100",
					Area:         "LOCAL",
					Address:      "cold78_status.VentSK_status",
				},
			},
			wantGuids: []string{
				"01749246-9617-585e-9e19-157ccad61ee2",
			},
			wantErr: nil,
		},
		{
			name:      "BAD#1",
			args:      args{dir: "testDirectory", file: "BAD1"},
			wantTsv:   nil,
			wantGuids: nil,
			wantErr:   constants.ErrNotTSV,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := createTempDir(tt.args.dir, t)
			if err != nil {
				t.Errorf("not creating temp dir: %v", err)
				return
			}

			defer func(dir string) {
				err := removeTempDir(dir)
				if err != nil {
					t.Errorf("not delete temp dir: %v", err)
					return
				}
			}(tempDir)

			file, err := createTempFile(tempDir, tt.args.file)
			if err != nil {
				t.Errorf("not creating temp file: %v", err)
				return
			}

			err = writeDataToFile(file, tt.wantTsv)
			if err != nil {
				t.Errorf("not writing temp file: %v", err)
				return
			}
			s := &Parser{
				dirFrom: config.Config{DirectoryFrom: tempDir}.DirectoryFrom,
			}
			tsvChan, guidChan, errChan := s.ParseFileAsync(tt.args.file)

			var gotTsv []shema.Tsv
			var gotGuids []string
			var gotErr error
			var wg sync.WaitGroup
			wg.Add(3) // Increase the counter to 3

			go func() {
				defer wg.Done()
				for tsv := range tsvChan {
					gotTsv = append(gotTsv, tsv)
				}
			}()

			go func() {
				defer wg.Done()
				for guid := range guidChan {
					gotGuids = append(gotGuids, guid)
				}
			}()

			// Add a new goroutine to read from the error channel
			go func() {
				defer wg.Done()
				for err := range errChan {
					if err != nil {
						gotErr = err
						break
					}
				}
			}()

			wg.Wait()

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("ParseFileAsync() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTsv, tt.wantTsv) {
				t.Errorf("ParseFileAsync() got = %v, want %v", gotTsv, tt.wantTsv)
			}
			if !reflect.DeepEqual(gotGuids, tt.wantGuids) {
				t.Errorf("ParseFileAsync() got1 = %v, want %v", gotGuids, tt.wantGuids)
			}
		})
	}
}

func createTempDir(dir string, t *testing.T) (string, error) {
	tempDir, err := os.MkdirTemp(".", dir)
	if err != nil {
		t.Errorf("not created directory")
		return "", err
	}
	return tempDir, nil
}

func createTempFile(dir, name string) (*os.File, error) {
	tempFile, err := os.Create(dir + "/" + name)
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}

func writeDataToFile(file *os.File, data []shema.Tsv) error {
	writer := csv.NewWriter(file)
	writer.Comma = '\t'

	for _, d := range data {
		record := []string{d.Number, d.MQTT, d.InventoryID, d.UnitGUID, d.MessageID, d.MessageText, d.Context,
			d.MessageClass, d.Level, d.Area, d.Address, d.Block, d.Type, d.Bit, d.InvertBit}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func removeTempDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}
