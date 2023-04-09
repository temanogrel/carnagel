package loadbalancer

import (
	"sort"
	"sync"
	"sync/atomic"

	"git.misc.vee.bz/carnagel/minerva/pkg"
	"git.misc.vee.bz/carnagel/minerva/pkg/strategy"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type loadBalancerService struct {
	app *minerva.Application
	log logrus.FieldLogger

	//different balancing strategies
	leastTotalBandwidthStrategy minerva.ServerLoadBalanceStrategy
	leastBandwidthStrategy      minerva.ServerLoadBalanceStrategy
	smallFileStrategy           minerva.ServerLoadBalanceStrategy
}

func NewService(app *minerva.Application) minerva.LoadBalancer {
	lbs := &loadBalancerService{
		app: app,
		log: app.Logger.WithField("component", "loadbalancer"),
	}

	lbs.leastTotalBandwidthStrategy = strategy.NewBandwidthStrategy(app.ServerCollection)
	lbs.leastBandwidthStrategy = strategy.NewBandwidthStrategy(app.ServerCollection)
	lbs.smallFileStrategy = lbs.leastTotalBandwidthStrategy

	return lbs
}

func (service *loadBalancerService) RecommendStorage(source minerva.Hostname, size uint64) (*minerva.Server, string, error) {

	minimumDiskSpace := service.getMinimumRequiredSpace()
	respectServersAvailability := service.getRespectServerAvailableForLoadBalancing()
	largeFileThreshold := service.getLargeFileThreshold()

	filters := []minerva.ServerFilterFunc{
		minerva.ServerOnlineInConsulFilter(),
		minerva.ServerMinimumFreeSpaceFilter(minimumDiskSpace),
	}

	// figure out which interface we should use
	internalOrExternalStrategy := strategy.NewInternalOrExternalInterfaceStrategy(service.app.ServerCollection)
	_, selectedInterface := internalOrExternalStrategy.Balance(source)

	var server *minerva.Server
	var selectedBandwidthFunc minerva.ServerBandwidthFunc

	// This is not the best idea performance wise for future use and should probably be removed
	// but we currently need to move data as quickly as possible, thus routing should be based on current metrics
	//
	// The initial idea was to balance the transfers based upon the total bandwidth being consumed by the server
	// which is great for camtube.co, but right now we need to run in fast mode
	if selectedInterface == minerva.ExternalInterface {
		selectedBandwidthFunc = minerva.ServerExternalIncomingBandwidth
	} else {
		selectedBandwidthFunc = minerva.ServerInternalIncomingBandwidth
	}

	// If the file size is above the threshold balance for large files
	if size > largeFileThreshold {

		// we only ignore servers available for load balancing if it's a large file
		if respectServersAvailability {
			filters = append(filters, minerva.ServerAvailableForLoadBalancing())
		}

		server = service.leastTotalBandwidthStrategy.Balance(false, minerva.DecideOnLowerTotalBandwidth, selectedBandwidthFunc, filters)
	} else {
		server = service.smallFileStrategy.Balance(true, minerva.DecideOnLowerTotalBandwidth, selectedBandwidthFunc, filters)
	}

	if server == nil {
		service.log.WithFields(logrus.Fields{
			"interface":                  selectedInterface,
			"minimumDiskSpace":           minimumDiskSpace,
			"largeFileThreshold":         largeFileThreshold,
			"respectServersAvailability": respectServersAvailability,
		}).Warn("Failed to recommend server for storage")

		return nil, selectedInterface, minerva.ServerNotAvailableErr
	}

	service.log.WithFields(logrus.Fields{
		"externalHostname":           server.ExternalHostname,
		"internalHostname":           server.InternalHostname,
		"freeSpace":                  server.FreeSpace,
		"interface":                  selectedInterface,
		"minimumDiskSpace":           minimumDiskSpace,
		"largeFileThreshold":         largeFileThreshold,
		"respectServersAvailability": respectServersAvailability,
	}).Debug("Recommended server for storage")

	return server, selectedInterface, nil
}

func (service *loadBalancerService) RecommendDownload(source minerva.Hostname, fileUUID string) (minerva.DownloadPath, error) {
	// lookup server which stores the file
	file, err := service.app.FileRepository.GetByUuid(uuid.FromStringOrNil(fileUUID))
	if err != nil {
		return minerva.DownloadPath{}, minerva.FileNotFoundErr
	}

	// figure out which interface we should use
	internalOrExternalStrategy := strategy.NewInternalOrExternalInterfaceStrategy(service.app.ServerCollection)

	foundOrigin, selectedInterface := internalOrExternalStrategy.Balance(file.Hostname)
	if foundOrigin == nil {
		return minerva.DownloadPath{}, minerva.FailedToLocateOriginServerErr
	}

	filters := []minerva.ServerFilterFunc{
		minerva.ServerOnlineInConsulFilter(),
	}

	var edge *minerva.Server

	// choose the least utilized server as an edge if the network is external
	if selectedInterface == minerva.ExternalInterface {
		edge = service.leastBandwidthStrategy.Balance(true, minerva.DecideOnLowerTotalBandwidth, minerva.ServersExternalBandwidth, filters)
	}

	downloadPath := minerva.DownloadPath{
		File:           file,
		Edge:           edge,
		Origin:         foundOrigin,
		InterfaceToUse: selectedInterface,
	}

	return downloadPath, nil
}

func (service *loadBalancerService) RedistributeData(top, bottom int, amountPerServer uint64) minerva.RedistributionReport {

	servers := service.getServersWithMostAvailableSpace()

	if bottom == 0 || bottom+top >= len(servers) {
		bottom = len(servers) - top
	}

	targets := servers[0:top]
	sources := servers[len(servers)-bottom:]

	report := minerva.RedistributionReport{
		TopServers:     targets,
		BottomServers:  sources,
		Redistribution: make(map[minerva.Hostname]*uint64),
	}

	// populate redistribution
	for _, server := range targets {
		report.Redistribution[server.ExternalHostname] = new(uint64)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(sources))

	for _, s := range sources {
		go func(server *minerva.Server) {
			defer wg.Done()

			files, err := service.app.FileRepository.GetOldestUpdatedAtByHostnameAndAccumulatedSize(server.ExternalHostname, amountPerServer)
			if err != nil {
				service.app.Logger.WithError(err).Error("Failed to retrieve files to be redistributed")
				return
			}

			var size uint64
			var targetIndex int

			for _, f := range files {
				size += f.Size
				target := targets[targetIndex]

				// Increment the redistribution to the server
				atomic.AddUint64(report.Redistribution[target.ExternalHostname], f.Size)

				// Reset target index once we loop through it
				if targetIndex == len(targets)-1 {
					targetIndex = 0
				}

				targetIndex++

				server.RelocateRequests <- minerva.RelocateRequest{
					Uuid:       f.Uuid,
					TargetHost: target.InternalHostname,
				}
			}

			service.app.Logger.WithFields(logrus.Fields{
				"hostname":   server.ExternalHostname,
				"size":       size,
				"targetSize": amountPerServer,
				"top":        top,
			}).Debug("Transferring files the most empty servers available")
		}(s)
	}

	wg.Wait()

	return report
}

func (service *loadBalancerService) getServersWithMostAvailableSpace() minerva.ServersSortedByFreeSpace {

	service.app.ServerCollection.ServersMtx.RLock()
	defer service.app.ServerCollection.ServersMtx.RUnlock()

	servers := make(minerva.ServersSortedByFreeSpace, 0, len(service.app.ServerCollection.Servers))

	for _, server := range service.app.ServerCollection.Servers {
		servers = append(servers, server)
	}

	sort.Sort(servers)

	return servers
}

func (service *loadBalancerService) getMinimumRequiredSpace() uint64 {
	// Get minimum free disk space with a default value of 50 if it fails
	return service.app.Consul.GetUint64("minerva/loadbalancer/minimum-free-diskspace", 50*1024*1024*1024)
}

func (service *loadBalancerService) getLargeFileThreshold() uint64 {
	// Get large file threshold with a default value of 1 mb
	return service.app.Consul.GetUint64("minerva/loadbalancer/large-file-threshold", 1024*1024)
}

func (service *loadBalancerService) getRespectServerAvailableForLoadBalancing() bool {
	return service.app.Consul.GetBool("minerva/loadbalancer/respect-server-availability", true)
}
