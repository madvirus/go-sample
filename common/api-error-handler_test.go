package common

import (
	"bytes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorHandlerTestSuite))
}

type ErrorHandlerTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *ErrorHandlerTestSuite) SetupTest() {
	router := gin.Default()
	router.POST("/invalidJsonBody", invalidJsonHandler)
	router.GET("/validationError", validationErrorHandler)
	router.GET("/noData", noDataHandler)
	suite.router = router
}

type JsonBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func invalidJsonHandler(c *gin.Context) {
	var req JsonBody
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": 10})
}

func validationErrorHandler(c *gin.Context) {
	errors := []ErrorField{
		{
			Name:    "name",
			Message: "messagename",
		},
		{
			Name:    "email",
			Message: "messageemail",
		},
	}
	ve := CreateValidationError(errors)
	HandleErrorResponse(c, ve)
}

func noDataHandler(c *gin.Context) {
	err := CreateNoDataFoundError("no data")
	HandleErrorResponse(c, err)
}

func (suite *ErrorHandlerTestSuite) TestInvalidJson() {
	w := httptest.NewRecorder()
	invalidJson := `{"name": "name", }`
	req, _ := http.NewRequest(http.MethodPost, "/invalidJsonBody", bytes.NewBuffer([]byte(invalidJson)))
	suite.router.ServeHTTP(w, req)

	log.Infof("response body: %s", w.Body.String())
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ErrorHandlerTestSuite) TestValidationError() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/validationError", nil)
	suite.router.ServeHTTP(w, req)

	log.Infof("response body: %s", w.Body.String())
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ErrorHandlerTestSuite) TestNoDataFoundError() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/noData", nil)
	suite.router.ServeHTTP(w, req)

	log.Infof("response body: %s", w.Body.String())
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}
