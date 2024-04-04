package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var UnmarshalTypeError *json.UnmarshalTypeError

func HandlerErr(c *gin.Context, err error) {
	if err != nil {
		switch {
		case errors.As(err, &UnmarshalTypeError):
			err := fmt.Sprintf("bad json %s", err)
			c.JSON(http.StatusBadRequest, err)
		default:
			c.JSON(http.StatusBadRequest, err)
		}
		return

	}

	c.Status(http.StatusOK)
	return
}
