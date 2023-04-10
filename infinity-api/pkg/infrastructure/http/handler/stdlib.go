package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type ApiResponse struct {
	Code    uint16            `json:"code"`
	Message string            `json:"message"`
	Error   error             `json:"error,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type ResponseMeta struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type CriteriaMeta struct {
	Limit      uint64
	LastSeenId uint64
}

var (
	RESP_SERVER_ERROR = &ApiResponse{
		Code:    500,
		Message: "Internal server error occured",
	}

	RESP_UNAUTHORIZED = &ApiResponse{
		Code:    401,
		Message: "Unauthorized",
	}

	RESP_FORBIDDEN = &ApiResponse{
		Code:    http.StatusForbidden,
		Message: "Unauthorized",
	}

	RESP_OBJECT_NOT_FOUND = &ApiResponse{
		Code:    404,
		Message: "Object not found",
	}

	RESP_MALFORMED_DATA = &ApiResponse{
		Code:    400,
		Message: "Data did not match expected structure",
	}

	RESP_MALFORMED_ID = &ApiResponse{
		Code:    400,
		Message: "Identifier was malformed",
	}
)

func jsonUnauthorizedResponse(w http.ResponseWriter) {
	jsonResponse(w, http.StatusUnauthorized, RESP_UNAUTHORIZED)
}

func jsonForbiddenResponse(w http.ResponseWriter) {
	jsonResponse(w, http.StatusForbidden, RESP_FORBIDDEN)
}

func jsonResponse(w http.ResponseWriter, status int, payload interface{}) error {

	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode json response"))

		return errors.Wrap(err, "Failed to encode json payload")
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	return nil
}
