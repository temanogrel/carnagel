package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

func ToRawMessage(data interface{}) (json.RawMessage, error) {
	resp, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal JSON")
	}

	return json.RawMessage(resp), nil
}

func CreateErrorResponse(message string, args ...interface{}) []byte {
	return CreateBroadcastResponse("error", map[string]string{
		"message": fmt.Sprintf(message, args...),
	})
}

func CreateBroadcastResponse(name string, payload interface{}) []byte {
	response, err := json.Marshal(Message{
		Name:    name,
		Payload: payload,
	})

	if err != nil {
		panic(err)
	}

	return response
}

func CreateRpcResponse(id uint32, payload json.RawMessage, success bool) []byte {
	response, err := json.Marshal(ProcedureResponse{
		Id:      id,
		Payload: payload,
		Success: success,
	})

	//todo: no panic
	if err != nil {
		panic(err)
	}

	return response
}
