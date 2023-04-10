package rpc

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	procedureResponseTime = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "websocket_rpc_duration_seconds",
			Help: "Number of seconds spent on each procedure call",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{"procedure"},
	)

	procedureErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_rpc_error_counter",
			Help: "Number of errors that have occurred over time",
		},
		[]string{"procedure"},
	)
)

func init() {
	prometheus.MustRegister(procedureResponseTime)
	prometheus.MustRegister(procedureErrorCount)
}

type server struct {
	logger       logrus.FieldLogger
	procedureMap ProcedureMap
}

func NewRpcServer(logger logrus.FieldLogger, procedureMap ProcedureMap) *server {
	s := &server{
		logger:       logger,
		procedureMap: procedureMap,
	}

	return s
}

func (s *server) ServeRequest(socket Socket, request *Request) {

	// measure time for prometheus
	startedAt := time.Now()

	procedure, ok := s.procedureMap[request.Name]
	if !ok {
		socket.Send(CreateErrorResponse("Invalid RPC '%s'", request.Name))
		return
	}

	if resp, err := procedure(request, socket); err == nil {
		socket.Send(CreateRpcResponse(request.Id, resp, true))
	} else {

		// Increment error count
		procedureErrorCount.
			WithLabelValues(request.Name).
			Inc()

		payload := FailedProcedurePayload{
			Error: err.Error(),
		}

		if message, err := ToRawMessage(payload); err == nil {
			socket.Send(CreateRpcResponse(request.Id, message, false))
		} else {

			s.logger.
				WithError(err).
				Error("Failed to convert to RawMessage")
		}
	}

	procedureResponseTime.
		WithLabelValues(request.Name).
		Observe(time.Since(startedAt).Seconds())
}
