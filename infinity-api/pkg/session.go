package infinity

import (
	"git.misc.vee.bz/carnagel/infinity-api/pkg/infrastructure/http/rpc"
	"github.com/sasha-s/go-deadlock"
	"github.com/satori/go.uuid"
)

type ActiveSessionCollection struct {
	Mtx      deadlock.RWMutex
	Sessions map[uuid.UUID]*rpc.Socket
}

func (c *ActiveSessionCollection) SendToUser(user uuid.UUID, message []byte) {
	c.Mtx.RLock()
	defer c.Mtx.RUnlock()

	if socket, ok := c.Sessions[user]; ok {
		go (*socket).Send(message)
	}
}
