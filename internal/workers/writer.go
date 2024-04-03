package workers

import (
	"fmt"
	"github.com/signintech/gopdf"
	"goTSVParser/config"
	"goTSVParser/internal/shema"
	"os"
	"path/filepath"
	"strings"
)

type Writer struct {
	dirTo   string
	dirFrom string
}

func NewWriter(cfg config.Config) *Writer {
	return &Writer{dirTo: cfg.DirectoryTo, dirFrom: cfg.DirectoryFrom}
}

func (s *Writer) WritePDF(tsv []shema.Tsv, unitGuid []string, filePath string) error {
	for _, guid := range unitGuid {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
		pdf.AddPage()

		defer pdf.Close()

		err := pdf.AddTTFFont("LiberationSerif-Regular", "resources/LiberationSerif-Regular.ttf")
		if err != nil {
			return fmt.Errorf("can't add font: %w", err)
		}

		err = pdf.SetFont("LiberationSerif-Regular", "", 14)
		if err != nil {
			return fmt.Errorf("can't set font: %w", err)
		}

		for _, t := range tsv {
			var resultArray []string
			if guid == t.UnitGUID {
				pdf.AddPage()

				resultArray = append(resultArray, "n: "+strings.TrimSpace(t.Number))
				resultArray = append(resultArray, "mqtt: "+strings.TrimSpace(t.MQTT))
				resultArray = append(resultArray, "invid: "+strings.TrimSpace(t.InventoryID))
				resultArray = append(resultArray, "unit_guid: "+strings.TrimSpace(t.UnitGUID))
				resultArray = append(resultArray, "msg_id: "+strings.TrimSpace(t.MessageID))
				resultArray = append(resultArray, "text: "+strings.TrimSpace(t.MessageText))
				resultArray = append(resultArray, "context: "+strings.TrimSpace(t.Context))
				resultArray = append(resultArray, "class: "+strings.TrimSpace(t.MessageClass))
				resultArray = append(resultArray, "level: "+strings.TrimSpace(t.Level))
				resultArray = append(resultArray, "area: "+strings.TrimSpace(t.Area))
				resultArray = append(resultArray, "addr: "+strings.TrimSpace(t.Address))
				resultArray = append(resultArray, "block: "+strings.TrimSpace(t.Block))
				resultArray = append(resultArray, "type: "+strings.TrimSpace(t.Type))
				resultArray = append(resultArray, "bit: "+strings.TrimSpace(t.Bit))
				resultArray = append(resultArray, "invert_bit: "+strings.TrimSpace(t.InvertBit))

				y := 20
				for _, str := range resultArray {
					pdf.SetXY(10, float64(y))
					err := pdf.Text(str)
					if err != nil {
						return fmt.Errorf("can't write string to PDF: %w", err)
					}
					y += 20
				}
			}
		}

		dir := filepath.Dir(strings.TrimPrefix(filePath, s.dirFrom))
		if err := os.MkdirAll(s.dirTo+dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		resultFile := s.dirTo + dir + "/" + guid + ".pdf"
		err = pdf.WritePdf(resultFile)
		if err != nil {
			return fmt.Errorf("failed to write result")
		}

	}
	return nil
}
