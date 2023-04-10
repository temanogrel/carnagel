package websocket

import (
	"bytes"
	"encoding/json"
	"github.com/sasha-s/go-deadlock"
	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/rpc"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the socketHandler.
type WrappedWebSocket struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// is the socket has been closed
	isClosed bool

	mtx deadlock.RWMutex

	// Values have the same purpose as context values
	values map[string]interface{}
}

func (s *WrappedWebSocket) Value(key string) interface{} {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if val, ok := s.values[key]; ok {
		return val
	}

	return nil
}

func (s *WrappedWebSocket) SetValue(key string, value interface{}) {
	s.mtx.Lock()
	s.values[key] = value
	s.mtx.Unlock()
}

func (s *WrappedWebSocket) Send(data []byte) {
	if !s.isClosed {
		s.send <- data
	}
}

func (s *WrappedWebSocket) Close() {
	s.isClosed = true

	close(s.send)
}

// readPump pumps messages from the websocket connection to the socketHandler.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (s *WrappedWebSocket) ReadPump(rpcServer rpc.RequestServer) {
	defer func() {
		s.conn.Close()
	}()

	s.conn.SetReadLimit(maxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(pongWait))

	s.conn.SetPongHandler(func(string) error {
		s.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		request := &rpc.Request{}
		if err := json.Unmarshal(message, request); err != nil {
			continue
		}

		go rpcServer.ServeRequest(s, request)
	}
}

// writePump pumps messages from the socketHandler to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (s *WrappedWebSocket) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		s.conn.Close()
	}()

	for {
		select {
		case message, ok := <-s.send:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The socketHandler closed the channel.
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := s.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
