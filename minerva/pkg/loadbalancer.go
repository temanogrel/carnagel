package minerva

import "github.com/pkg/errors"

var (
	FailedToLocateOriginServerErr = errors.New("Failed to locate the origin server of the download request")
)

const (
	InternalInterface = "internal"
	ExternalInterface = "external"
)

type DownloadPath struct {
	File           *File
	Edge           *Server
	Origin         *Server
	InterfaceToUse string
}

type RedistributionReport struct {
	TopServers     ServersSortedByFreeSpace `json:"topServers"`
	BottomServers  ServersSortedByFreeSpace `json:"bottomServers"`
	Redistribution map[Hostname]*uint64     `json:"redistribution"`
}

type LoadBalancer interface {
	// RecommendStorage recommends the best server to use based on current BW criteria
	// The return values are the target server, interface and an error if one occurred
	RecommendStorage(source Hostname, size uint64) (*Server, string, error)

	// RecommendDownload recommends the best path a request can take to maximize bandwidth usage
	// of the origin server
	RecommendDownload(source Hostname, fileUUID string) (DownloadPath, error)

	// RedistributeData will redistribute files between the top x servers for the amount per server specified
	RedistributeData(top, bottom int, amountPerServer uint64) RedistributionReport
}

type ServerLoadBalanceStrategy interface {
	Balance(keepServerAvailable bool, decide LoadBalancerDecisionFunc, determineServerBandWidth ServerBandwidthFunc, filters []ServerFilterFunc) *Server
}

type LoadBalancerDecisionFunc func(currentBest, candidate *Server, determineServerBandwidth ServerBandwidthFunc) *Server

func DecideOnLowerTotalBandwidth(currentBest, candidate *Server, determineServerBandwidth ServerBandwidthFunc) *Server {
	//if bestServers used bandwidth is higher than servers, prefer server, or set any if not already done
	if determineServerBandwidth(candidate) < determineServerBandwidth(currentBest) {
		return candidate
	}

	return currentBest
}
