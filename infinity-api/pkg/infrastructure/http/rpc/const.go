package rpc

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var (
	RequestInvalidBodyErr = errors.New("Invalid request body provided")
)

type Procedure func(*Request, Socket) (json.RawMessage, error)
type ProcedureMap map[string]Procedure

type SocketConnectedCallback func(socket Socket)
type SocketDisconnectedCallback func(socket Socket)

type Socket interface {
	Send([]byte)
	Close()

	Value(string) interface{}
	SetValue(string, interface{})
}

type RequestServer interface {
	ServeRequest(Socket, *Request)
}

type Request struct {
	Id      uint32          `json:"id"`
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

type ProcedureResponse struct {
	Id      uint32          `json:"id"`
	Success bool            `json:"success"`
	Payload json.RawMessage `json:"payload"`
}

type FailedProcedurePayload struct {
	Error string `json:"error"`
}

type Message struct {
	Name    string      `json:"name"`
	Payload interface{} `json:"payload"`
}
