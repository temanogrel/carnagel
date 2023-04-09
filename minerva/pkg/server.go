package minerva

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var ServerNotAvailableErr = errors.New("No server is currently available for storage")

const MAX_PENDING_OPERATIONS = 100000

type Hostname string

type RelocateRequest struct {
	Uuid       uuid.UUID
	TargetHost Hostname
}

type Server struct {
	FreeSpace uint64 `json:"freeSpace"`

	InternalIp            string   `json:"internalIp"`
	InternalHostname      Hostname `json:"internalHostname"`
	InternalBandwidthUp   uint64   `json:"internalBandwidthUp"`
	InternalBandwidthDown uint64   `json:"internalBandwidthDown"`

	ExternalIp            string   `json:"externalIp"`
	ExternalHostname      Hostname `json:"externalHostname"`
	ExternalBandwidthUp   uint64   `json:"externalBandwidthUp"`
	ExternalBandwidthDown uint64   `json:"externalBandwidthDown"`

	OnlineInConsul            bool `json:"onlineInConsul"`
	AvailableForLoadBalancing bool `json:"availableForLoadBalancing"`

	UpdatedAt time.Time `json:"updatedAt"`

	UploadRequests   chan uuid.UUID       `json:"-"`
	DeletionRequests chan uuid.UUID       `json:"-"`
	RelocateRequests chan RelocateRequest `json:"-"`
}

func NewServer(internalHostname, externalHostname Hostname) *Server {
	return &Server{
		InternalHostname: internalHostname,
		ExternalHostname: externalHostname,

		UpdatedAt: time.Now(),

		UploadRequests:   make(chan uuid.UUID, MAX_PENDING_OPERATIONS),
		DeletionRequests: make(chan uuid.UUID, MAX_PENDING_OPERATIONS),
		RelocateRequests: make(chan RelocateRequest, MAX_PENDING_OPERATIONS),
	}
}

func (s *Server) TotalBandwidth() uint64 {
	return s.TotalExternalBandwidth() + s.TotalInternalBandwidth()
}

func (s *Server) TotalInternalBandwidth() uint64 {
	return s.InternalBandwidthDown + s.InternalBandwidthUp
}

func (s *Server) TotalExternalBandwidth() uint64 {
	return s.ExternalBandwidthDown + s.ExternalBandwidthUp
}

type ServerBandwidthFunc func(server *Server) uint64

func ServersTotalBandwidth(server *Server) uint64 {
	return server.TotalBandwidth()
}

func ServersInternalBandwidth(server *Server) uint64 {
	return server.TotalBandwidth()
}

func ServersExternalBandwidth(server *Server) uint64 {
	return server.TotalBandwidth()
}

func ServerExternalIncomingBandwidth(server *Server) uint64 {
	return server.ExternalBandwidthDown
}

func ServerInternalIncomingBandwidth(server *Server) uint64 {
	return server.InternalBandwidthDown
}

type ServerFilterFunc func(*Server) bool

func ServerMinimumFreeSpaceFilter(minimumSpace uint64) ServerFilterFunc {
	return func(s *Server) bool {
		if s.FreeSpace >= minimumSpace {
			return true
		}

		return false
	}
}

func ServerOnlineInConsulFilter() ServerFilterFunc {
	return func(s *Server) bool {
		return s.OnlineInConsul
	}
}

func ServerAvailableForLoadBalancing() ServerFilterFunc {
	return func(s *Server) bool {
		return s.AvailableForLoadBalancing
	}
}

type ServerDiscovery interface {
	Run(ctx context.Context)
}

type ServerMetricCollector interface {
	Run(ctx context.Context)
}

type Servers struct {
	InternalIpMap map[string]*Server
	ExternalIpMap map[string]*Server
	Servers       map[Hostname]*Server
	ServersMtx    sync.RWMutex
}

func NewServers() *Servers {
	s := &Servers{
		Servers:       make(map[Hostname]*Server),
		InternalIpMap: make(map[string]*Server),
		ExternalIpMap: make(map[string]*Server),
	}

	return s
}

type ServersSortedByFreeSpace []*Server

func (p ServersSortedByFreeSpace) Len() int           { return len(p) }
func (p ServersSortedByFreeSpace) Less(i, j int) bool { return p[i].FreeSpace > p[j].FreeSpace }
func (p ServersSortedByFreeSpace) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
