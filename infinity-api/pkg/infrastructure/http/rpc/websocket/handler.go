package websocket

import (
	"net/http"

	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/rpc"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// socketHandler maintains the set of active sockets and broadcasts messages to the
type SocketHandler struct {
	logger logrus.FieldLogger
	rpc    rpc.RequestServer

	onSocketConnected    rpc.SocketConnectedCallback
	onSocketDisconnected rpc.SocketDisconnectedCallback
}

func NewSocketHandler(
	logger logrus.FieldLogger,
	rpcServer rpc.RequestServer,
	onSocketConnected rpc.SocketConnectedCallback,
	onSocketDisconnected rpc.SocketDisconnectedCallback,
) *SocketHandler {
	return &SocketHandler{
		rpc:    rpcServer,
		logger: logger,

		onSocketConnected:    onSocketConnected,
		onSocketDisconnected: onSocketDisconnected,
	}
}

func (h *SocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.WithError(err).Warn("Failed to upgrade request for websocket connection")
		return
	}

	socket := &WrappedWebSocket{
		conn:   conn,
		send:   make(chan []byte, 256),
		values: make(map[string]interface{}),
	}

	h.onSocketConnected(socket)

	go socket.WritePump()
	socket.ReadPump(h.rpc)

	h.onSocketDisconnected(socket)
}
