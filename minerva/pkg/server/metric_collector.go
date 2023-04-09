package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"git.misc.vee.bz/carnagel/minerva/pkg"
	"git.misc.vee.bz/carnagel/minerva/pkg/internal"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

type metricCollector struct {
	app    *minerva.Application
	logger logrus.FieldLogger
}

func NewMetricCollector(app *minerva.Application) minerva.ServerMetricCollector {
	return &metricCollector{
		app:    app,
		logger: app.Logger.WithField("component", "server_metric_collector"),
	}
}

func (collector *metricCollector) Run(ctx context.Context) {
	timer := time.NewTimer(5 * time.Second)

	// Don't wait 2 seconds for the first runt
	collector.poll(ctx)

	for {
		select {
		case <-ctx.Done():
			timer.Stop()

		case <-timer.C:
			collector.poll(ctx)

			timer.Reset(time.Second * 5)
		}
	}
}

func (collector *metricCollector) queryBandwidth(ctx context.Context, metric, device string) (model.Vector, error) {
	query := fmt.Sprintf("irate(%s{node=~'(edge|bs)[0-9]+', device=~'%s'}[5s])", metric, device)

	vector, err := collector.app.PrometheusClient.Query(ctx, query, time.Now())
	if err != nil {

		collector.logger.
			WithFields(logrus.Fields{
				"metric": metric,
				"device": device,
			}).
			WithError(err).
			Error("Failed to execute bandwidth query")

		return nil, errors.Wrap(err, "failed to query bandwidth")
	}

	return vector.(model.Vector), nil
}

func (collector *metricCollector) poll(ctx context.Context) {

	startedAt := time.Now()
	timeout := time.Millisecond * 1000

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(5)

	go collector.queryFreeDiskSpace(ctx, wg)
	go collector.queryIncomingInternalBandwidth(ctx, wg)
	go collector.queryIncomingInternalBandwidth(ctx, wg)
	go collector.queryOutgoingExternalBandwidth(ctx, wg)
	go collector.queryOutgoingExternalBandwidth(ctx, wg)

	wg.Wait()

	collector.app.ServerCollection.ServersMtx.Lock()

	// enable all servers for load balancing again
	for _, s := range collector.app.ServerCollection.Servers {
		s.AvailableForLoadBalancing = true
	}

	collector.app.ServerCollection.ServersMtx.Unlock()

	if ctx.Err() == nil {
		collector.app.Logger.
			WithField("component", "metric_collector").
			WithField("duration", time.Since(startedAt).Seconds()).
			Debug("Polled prometheus")

	} else {
		collector.app.Logger.
			WithField("component", "metric_collector").
			Warnf("Failed to poll prometheus, timed out after %d ms", timeout)
	}
}

func (collector *metricCollector) queryFreeDiskSpace(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	freeSpaceQuery := "sort_desc(node_filesystem_free{mountpoint='/data',node=~'(edge|bs)[0-9]+'})"
	freeSpaceVector, err := collector.app.PrometheusClient.Query(ctx, freeSpaceQuery, time.Now())
	if err != nil {
		collector.logger.WithError(err).Warn("Failed to execute filesystem query")
		return
	}

	collector.app.ServerCollection.ServersMtx.Lock()

	for _, vector := range freeSpaceVector.(model.Vector) {
		instance := internal.InstanceToIp(string(vector.Metric["instance"]))

		server, ok := collector.app.ServerCollection.InternalIpMap[instance]
		if ok {
			server.FreeSpace = uint64(vector.Value)
		}
	}

	collector.app.ServerCollection.ServersMtx.Unlock()
}

func (collector *metricCollector) queryIncomingExternalBandwidth(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// incoming on eth0
	vector, err := collector.queryBandwidth(ctx, "node_network_receive_bytes", "eth0")
	if err != nil {
		collector.logger.WithError(err).Warn("Failed to execute node_network_receive_bytes query")
		return
	}

	collector.app.ServerCollection.ServersMtx.Lock()

	for _, vector := range vector {
		instance := internal.InstanceToIp(string(vector.Metric["instance"]))

		server, ok := collector.app.ServerCollection.InternalIpMap[instance]
		if ok {
			server.ExternalBandwidthDown = uint64(vector.Value)
		}
	}

	collector.app.ServerCollection.ServersMtx.Unlock()
}

func (collector *metricCollector) queryOutgoingExternalBandwidth(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// outgoing on eth0
	vector, err := collector.queryBandwidth(ctx, "node_network_transmit_bytes", "eth0")
	if err != nil {
		collector.logger.WithError(err).Warn("Failed to execute node_network_transmit_bytes query")
		return
	}

	collector.app.ServerCollection.ServersMtx.Lock()

	for _, vector := range vector {
		instance := internal.InstanceToIp(string(vector.Metric["instance"]))

		server, ok := collector.app.ServerCollection.InternalIpMap[instance]
		if ok {
			server.ExternalBandwidthUp = uint64(vector.Value)
		}
	}

	collector.app.ServerCollection.ServersMtx.Unlock()
}

func (collector *metricCollector) queryIncomingInternalBandwidth(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// incoming on eth1..3
	vector, err := collector.queryBandwidth(ctx, "node_network_receive_bytes", "eth[1-3]")
	if err != nil {
		collector.logger.WithError(err).Warn("Failed to execute node_network_receive_bytes query")
		return
	}

	collector.app.ServerCollection.ServersMtx.Lock()

	for _, vector := range vector {
		instance := internal.InstanceToIp(string(vector.Metric["instance"]))

		server, ok := collector.app.ServerCollection.InternalIpMap[instance]
		if ok {
			server.InternalBandwidthDown = uint64(vector.Value)
		}
	}

	collector.app.ServerCollection.ServersMtx.Unlock()
}

func (collector *metricCollector) queryOutgoingInternalBandwidth(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// outgoing on eth1..3
	vector, err := collector.queryBandwidth(ctx, "node_network_transmit_bytes", "eth[1-3]")
	if err != nil {
		collector.logger.WithError(err).Warn("Failed to execute node_network_transmit_bytes query")
		return
	}

	collector.app.ServerCollection.ServersMtx.Lock()

	for _, vector := range vector {
		instance := internal.InstanceToIp(string(vector.Metric["instance"]))

		server, ok := collector.app.ServerCollection.InternalIpMap[instance]
		if ok {
			server.InternalBandwidthUp = uint64(vector.Value)
		}
	}

	collector.app.ServerCollection.ServersMtx.Unlock()
}
