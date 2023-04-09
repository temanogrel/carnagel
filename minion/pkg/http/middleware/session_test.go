package middleware

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/suite"
	"gopkg.in/stretchr/testify.v1/assert"
)

type sessionMiddlewareTest struct {
	suite.Suite

	middleware *sessionMiddleware
}

func (suite *sessionMiddlewareTest) BeforeTest(suiteName, testName string) {
	suite.middleware = &sessionMiddleware{}
}

func (suite *sessionMiddlewareTest) TestWithMissingSession() {
	request := httptest.NewRequest("GET", "http://example.com/foo", nil)
	response := httptest.NewRecorder()

	suite.middleware.ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusBadRequest, response.Code)
	assert.Equal(suite.T(), "Missing session", response.Body.String())
}

func (suite *sessionMiddlewareTest) TestWithInvalidSession() {
	request := httptest.NewRequest("GET", "http://example.com/foo", nil)
	request.URL.RawQuery = "session=invalid-session"
	response := httptest.NewRecorder()

	suite.middleware.ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusUnauthorized, response.Code)
	assert.Equal(suite.T(), "Invalid session provided", response.Body.String())
}

func TestSessionMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(sessionMiddlewareTest))
}
