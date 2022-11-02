package common

import (
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// HandleErrorResponse 에러 타입에 따라 알맞은 결과를 응답하는 gin 에러 핸들러
func HandleErrorResponse(c *gin.Context, err error) {
	log.Errorf("error occured: %T", err)
	switch e := err.(type) {
	case *json.SyntaxError:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": e.Error(),
		})
	case *ValidationError:
		c.JSON(http.StatusBadRequest, e)
	case *NoDataFoundError:
		c.JSON(http.StatusNotFound, e)
	default:
		sentry.CaptureException(e)
		c.JSON(http.StatusInternalServerError, e)
	}
}
