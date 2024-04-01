package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"goTSVParser/config"
	"goTSVParser/internal/domains/mocks"
	"goTSVParser/internal/shema"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type serviceMock func(c *mocks.Service)

func TestHandler_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		body        shema.Request
		serviceMock serviceMock
		wantCode    int
		want        [][]shema.Tsv
	}{

		{
			name: "OK#1",
			body: shema.Request{
				UnitGUID: "ajsuiwp18203475nmgbdxgsk",
				Limit:    1,
				Page:     2,
			},
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetAll", mock.Anything, shema.Request{UnitGUID: "ajsuiwp18203475nmgbdxgsk", Limit: 1, Page: 2}).Return([][]shema.Tsv{{shema.Tsv{UnitGUID: "ajsuiwp18203475nmgbdxgsk"}}, {shema.Tsv{UnitGUID: "ajsuiwp18203475nmgbdxgsk"}}}, nil).Times(1)
			},
			wantCode: http.StatusOK,
			want:     [][]shema.Tsv{{shema.Tsv{UnitGUID: "ajsuiwp18203475nmgbdxgsk"}}, {shema.Tsv{UnitGUID: "ajsuiwp18203475nmgbdxgsk"}}},
		},
		{
			name: "OK#2",
			body: shema.Request{
				UnitGUID: "ajsuiwp18203475nmgbdxgsk",
				Limit:    1,
				Page:     3,
			},
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetAll", mock.Anything, shema.Request{UnitGUID: "ajsuiwp18203475nmgbdxgsk", Limit: 1, Page: 3}).Return([][]shema.Tsv{{shema.Tsv{UnitGUID: "ajsuiwp18203475nmgbdxgsk"}}}, nil).Times(1)
			},
			wantCode: http.StatusOK,
			want:     [][]shema.Tsv{{shema.Tsv{UnitGUID: "ajsuiwp18203475nmgbdxgsk"}}},
		},
		{
			name: "BAD#1",
			body: shema.Request{
				UnitGUID: "",
				Limit:    1,
				Page:     2,
			},
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetAll", mock.Anything, shema.Request{UnitGUID: "", Limit: 1, Page: 2}).Return(nil, errors.New("invalid data")).Times(1)
			},
			wantCode: http.StatusBadRequest,
			want:     nil,
		},
		{
			name: "BAD#2",
			body: shema.Request{
				UnitGUID: "1yua683",
				Limit:    1,
				Page:     1,
			},
			serviceMock: func(c *mocks.Service) {
				c.Mock.On("GetAll", mock.Anything, shema.Request{UnitGUID: "1yua683", Limit: 1, Page: 1}).Return(nil, errors.New("not found")).Times(1)
			},
			wantCode: http.StatusBadRequest,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			service := mocks.NewService(t)
			h := NewHandler(service, config.Config{})
			tt.serviceMock(service)

			path := "/api/all"
			g.POST(path, h.GetAll)
			b, err := json.Marshal(tt.body)
			if err != nil {
				fmt.Errorf("failed json")
				return
			}
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(string(b)))

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}

			wantResponse, err := json.Marshal(tt.want)
			if err != nil {
				fmt.Errorf("failed json")
				return
			}

			response, err := json.Marshal(w.Body)
			if err != nil {
				fmt.Errorf("failed json")
				return
			}

			if bytes.Equal(wantResponse, response) {
				t.Errorf("got %s, want %s", w.Body, tt.want)
			}
		})
	}
}
