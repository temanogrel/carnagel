package server

import (
	"context"
	"sync"
	"time"

	"git.misc.vee.bz/carnagel/minerva/pkg"
	consulapi "github.com/hashicorp/consul/api"
)

type discovery struct {
	app *minerva.Application

	pollTicker *time.Ticker
}

func NewServiceService(application *minerva.Application) minerva.ServerDiscovery {
	return &discovery{app: application}
}

func (discovery *discovery) Run(ctx context.Context) {

	timer := time.NewTimer(time.Second * 3)

	discovery.poll()

	for {
		select {
		case <-ctx.Done():
			return

		case <-timer.C:
			discovery.poll()

			timer.Reset(time.Second * 3)
		}
	}
}

func (discovery *discovery) poll() {

	services, _, err := discovery.app.Consul.API().Catalog().Service("minion", "", &consulapi.QueryOptions{})
	if err != nil {
		discovery.app.Logger.WithError(err).Errorf("Failed to retrieve services from consul")
		return
	}

	startedAt := time.Now()

	discovery.app.ServerCollection.ServersMtx.Lock()

	// Mark all servers as offline and have the next loop only enable the
	// servers that are still online
	for _, server := range discovery.app.ServerCollection.Servers {
		server.OnlineInConsul = false
	}

	wg := sync.WaitGroup{}

	for _, service := range services {

		internalHostname, ok := service.NodeMeta["internal_hostname"]
		if !ok {
			discovery.app.Logger.WithField("node", service.Node).Warn("Missing internal_hostname")
			continue
		}

		externalHostname, ok := service.NodeMeta["external_hostname"]
		if !ok {
			discovery.app.Logger.WithField("node", service.Node).Warn("Missing external_hostname")
			continue
		}

		if server, ok := discovery.app.ServerCollection.Servers[minerva.Hostname(externalHostname)]; ok {
			server.UpdatedAt = time.Now()
			server.OnlineInConsul = true
			continue
		}

		server := minerva.NewServer(minerva.Hostname(internalHostname), minerva.Hostname(externalHostname))
		server.InternalIp = service.TaggedAddresses["lan"]
		server.ExternalIp = service.TaggedAddresses["wan"]
		server.UpdatedAt = time.Now()
		server.OnlineInConsul = true
		server.AvailableForLoadBalancing = true

		discovery.app.ServerCollection.InternalIpMap[server.InternalIp] = server
		discovery.app.ServerCollection.ExternalIpMap[server.ExternalIp] = server
		discovery.app.ServerCollection.Servers[minerva.Hostname(externalHostname)] = server

		// Add all the pending files in a goroutine so we can run this a lot faster
		if !discovery.app.DevMode {
			go func() {
				wg.Add(1)
				defer wg.Done()

				// todo: run this operation in a goroutine so we can speed up bootstrapping
				files, err := discovery.app.FileRepository.GetWithPendingOperations(server.ExternalHostname, minerva.MAX_PENDING_OPERATIONS)
				if err != nil {
					discovery.app.Logger.
						WithError(err).
						Error("Failed to retrieve pending file operations")

					return
				}

				for _, file := range files {

					logger := discovery.app.Logger.WithField("file", file).WithField("server", server)

					if file.PendingDeletion {
						logger.Debug("Adding deletion request from the database")
						server.DeletionRequests <- file.Uuid
					}

					if file.PendingUpload {
						logger.Debug("Adding upload request from the database")
						server.UploadRequests <- file.Uuid
					}
				}
			}()
		}
	}

	wg.Wait()
	discovery.app.ServerCollection.ServersMtx.Unlock()
	discovery.app.Logger.
		WithField("duration", time.Since(startedAt).Seconds()).
		Debug("Polled consul for minions")
}
