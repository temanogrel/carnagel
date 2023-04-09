package strategy

import "git.misc.vee.bz/carnagel/minerva/pkg"

type internalOrExternalInterfaceStrategy struct {
	serverCollection *minerva.Servers
}

func NewInternalOrExternalInterfaceStrategy(servers *minerva.Servers) *internalOrExternalInterfaceStrategy {
	s := &internalOrExternalInterfaceStrategy{
		serverCollection: servers,
	}
	return s
}

func (s *internalOrExternalInterfaceStrategy) Balance(source minerva.Hostname) (*minerva.Server, string) {
	var foundOrigin *minerva.Server
	var selectedInterface string

	//see if the source inside our outside the storage network
	s.serverCollection.ServersMtx.RLock()
	server, ok := s.serverCollection.Servers[source]
	if !ok {
		//we could not find the hostname, so the request came from outside the network

		//find the server which has the requested external hostname
		for _, machine := range s.serverCollection.Servers {
			if machine.ExternalHostname == source {
				foundOrigin = machine
				break
			}
		}
		selectedInterface = "external"
	} else {
		foundOrigin = server
		selectedInterface = "internal"
	}
	s.serverCollection.ServersMtx.RUnlock()

	return foundOrigin, selectedInterface
}
