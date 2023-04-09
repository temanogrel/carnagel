package strategy

import "git.misc.vee.bz/carnagel/minerva/pkg"

type bandwidthStrategy struct {
	serverCollection *minerva.Servers
}

func NewBandwidthStrategy(servers *minerva.Servers) *bandwidthStrategy {
	s := &bandwidthStrategy{
		serverCollection: servers,
	}
	return s
}

func (s *bandwidthStrategy) Balance(
	keepServerAvailable bool,
	decide minerva.LoadBalancerDecisionFunc,
	determineServerBandWidth minerva.ServerBandwidthFunc,
	serverFilters []minerva.ServerFilterFunc) *minerva.Server {
	//get the server with the lowest Bandwidth usage
	var bestServer *minerva.Server

	filteredServers := make([]*minerva.Server, 0, 8)

	//filter out servers that match our "hard" criteria
	s.serverCollection.ServersMtx.RLock()

ServerLoop:
	for _, server := range s.serverCollection.Servers {
		for _, filter := range serverFilters {
			if filter(server) == false {
				continue ServerLoop
			}
		}

		filteredServers = append(filteredServers, server)
	}

	s.serverCollection.ServersMtx.RUnlock()

	//decide which one to use based on "soft" criteria
	s.serverCollection.ServersMtx.Lock()
	for _, server := range filteredServers {
		if bestServer == nil {
			bestServer = server
			continue
		}

		bestServer = decide(bestServer, server, determineServerBandWidth)
	}

	// mark server (if any) as unavailable for load balancer
	if bestServer != nil && keepServerAvailable == false {
		bestServer.AvailableForLoadBalancing = false
	}

	s.serverCollection.ServersMtx.Unlock()

	return bestServer
}
