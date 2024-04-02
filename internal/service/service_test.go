package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"goTSVParser/internal/constants"
	"goTSVParser/internal/domains/mocks"
	"goTSVParser/internal/shema"
	"reflect"
	"testing"
)

type storageMock[A any] func(c *mocks.Storage, args A)

func TestService_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		args        shema.Request
		storageMock storageMock[shema.Request]
		wantArr     [][]shema.Tsv
		wantErr     error
	}{
		{
			name: "OK1",
			args: shema.Request{
				UnitGUID: "01749246-9617-585e-9e19-157ccad61ee2",
				Limit:    1,
				Page:     1,
			},

			storageMock: func(c *mocks.Storage, r shema.Request) {
				c.Mock.On("GetAllGuids", mock.Anything, r.UnitGUID).Return([]shema.Tsv{
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
				}, nil).Times(1)
			},
			wantArr: [][]shema.Tsv{{
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
				}}},
			wantErr: nil,
		},
		{
			name: "OK2",
			args: shema.Request{
				UnitGUID: "01749246-9617-585e-9e19-157ccad61ee2",
				Limit:    1,
				Page:     2,
			},

			storageMock: func(c *mocks.Storage, r shema.Request) {
				c.Mock.On("GetAllGuids", mock.Anything, r.UnitGUID).Return([]shema.Tsv{
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
					{
						Number: "7",
					},
					{
						Number: "8",
					},
				}, nil).Times(1)
			},
			wantArr: [][]shema.Tsv{{
				{
					Number: "7",
				},
			}, {
				{
					Number: "8",
				},
			}},
			wantErr: nil,
		},
		{
			name: "BAD1",
			args: shema.Request{
				UnitGUID: "ahsyuiflgh",
				Limit:    1,
				Page:     1,
			},

			storageMock: func(c *mocks.Storage, r shema.Request) {
				c.Mock.On("GetAllGuids", mock.Anything, r.UnitGUID).Return(nil, errors.New("error getting")).Times(1)
			},
			wantArr: nil,
			wantErr: constants.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)
			tt.storageMock(storage, tt.args)
			logger, err := zap.NewProduction()

			service := Service{
				storage: storage,
				logger:  logger,
			}
			ctx := context.Background()
			tsvs, err := service.GetAll(ctx, tt.args)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tsvs, tt.wantArr) {
				t.Errorf("got %v, want %v", tsvs, tt.wantArr)
			}
		})
	}
}
