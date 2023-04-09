package http

import (
	"crypto/tls"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	*retryablehttp.Client

	logger logrus.FieldLogger
}

func NewHttpClient(logger logrus.FieldLogger) *HttpClient {
	client := &HttpClient{logger: logger}

	httpClient := retryablehttp.NewClient()
	// Self signed certificate so disable verification
	httpClient.HTTPClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	httpClient.ErrorHandler = client.logError

	client.Client = httpClient

	return client
}

func (c *HttpClient) LogUnexpectedResponseCode(request *http.Request, response *http.Response) {
	requestData := map[string]interface{}{}

	responseData := map[string]interface{}{
		"status":  response.StatusCode,
		"headers": response.Header,
	}

	if body, err := ioutil.ReadAll(response.Body); err == nil {
		responseData["body"] = body
	}

	if request != nil {
		requestData["url"] = request.URL.String()
		requestData["headers"] = request.Header
	}

	c.logger.WithFields(logrus.Fields{
		"request":  requestData,
		"response": responseData,
	}).Error("Unexpected response code received")
}

func (c *HttpClient) logError(response *http.Response, err error, numTries int) (*http.Response, error) {
	c.logger.WithError(err).Errorf("Failed to send http request")

	return response, err
}
