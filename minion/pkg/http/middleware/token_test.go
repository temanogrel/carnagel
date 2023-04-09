package middleware

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/suite"
	"gopkg.in/stretchr/testify.v1/assert"
)

type tokenMiddlewareTest struct {
	suite.Suite

	middleware *tokenMiddleware
}

func (suite *tokenMiddlewareTest) BeforeTest(suiteName, testName string) {
	suite.middleware = &tokenMiddleware{
		expectedToken: "validToken",
	}
}

func (suite *tokenMiddlewareTest) TestWithInvalidToken() {
	request := httptest.NewRequest("GET", "http://example.com/foo", nil)
	request.Header.Set("Authorization", "invalidToken")

	response := httptest.NewRecorder()

	suite.middleware.ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusUnauthorized, response.Code)
}

func TestTokenMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(tokenMiddlewareTest))
}
