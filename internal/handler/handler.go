package handler

import (
	"github.com/gin-gonic/gin"
	"goTSVParser/config"
	"goTSVParser/internal/domains"
	"goTSVParser/internal/shema"
	"net/http"
)

type Handler struct {
	service domains.Service
	engine  *gin.Engine
	config  config.Config
}

func NewHandler(service domains.Service, cnf config.Config) *Handler {
	router := gin.Default()
	h := &Handler{
		service: service,
		engine:  router,
		config:  cnf,
	}
	Route(router, h)
	return h
}

func (s *Handler) Start() error {
	err := s.engine.Run(s.config.Host)
	if err != nil {
		return err
	}
	return nil
}

func (s *Handler) GetAll(c *gin.Context) {
	var r shema.Request
	err := c.ShouldBindJSON(&r)
	if err != nil {
		HandlerErr(c, err)
		return
	}
	ctx := c.Request.Context()
	result, err := s.service.GetAll(ctx, r)
	if err != nil {
		HandlerErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)

}
