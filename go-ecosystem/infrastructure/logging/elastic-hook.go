package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)

var (
	// Fired if the
	// index is not created
	ErrCannotCreateIndex = fmt.Errorf("Cannot create index")
)

const FlushInterval = time.Second * 15

type IndexNameFunc func() string

// ElasticHook is a logrus
// hook for ElasticSearch
type ElasticHook struct {
	processor *elastic.BulkProcessor
	host      string
	index     IndexNameFunc
	levels    []logrus.Level
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewElasticHook(client *elastic.Client, host string, level logrus.Level, indexFunc IndexNameFunc) (*ElasticHook, error) {
	levels := []logrus.Level{}
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	} {
		if l <= level {
			levels = append(levels, l)
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(indexFunc()).Do(ctx)
	if err != nil {
		// Handle error
		return nil, err
	}

	if !exists {
		createIndex, err := client.CreateIndex(indexFunc()).Do(ctx)
		if err != nil {
			return nil, err
		}
		if !createIndex.Acknowledged {
			return nil, ErrCannotCreateIndex
		}
	}

	// from elastic docs: "If you want the bulk processor to
	// operate completely asynchronously, set both BulkActions and BulkSize to
	// -1 and set the FlushInterval to a meaningful interval."
	processor, err := client.BulkProcessor().
		Name(host).
		BulkActions(-1).
		BulkSize(-1).
		FlushInterval(FlushInterval).
		Do(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to create bulk processor")
	}

	return &ElasticHook{
		processor: processor,
		index:     indexFunc,
		host:      host,
		levels:    levels,
		ctx:       ctx,
		ctxCancel: cancel,
	}, nil
}

// Fire is required to implement
// Logrus hook
func (hook *ElasticHook) Fire(entry *logrus.Entry) error {
	level := entry.Level.String()

	if e, ok := entry.Data[logrus.ErrorKey]; ok && e != nil {
		if err, ok := e.(error); ok {
			entry.Data[logrus.ErrorKey] = err.Error()
		}
	}

	msg := struct {
		Host      string
		Timestamp string `json:"@timestamp"`
		Message   string
		Data      logrus.Fields
		Level     string
	}{
		hook.host,
		entry.Time.UTC().Format(time.RFC3339Nano),
		entry.Message,
		entry.Data,
		strings.ToUpper(level),
	}

	r := elastic.NewBulkIndexRequest().Index(hook.index()).Type("log").Doc(msg)
	hook.processor.Add(r)

	return nil
}

// Required for logrus
// hook implementation
func (hook *ElasticHook) Levels() []logrus.Level {
	return hook.levels
}

// Cancels all calls to
// elastic
func (hook *ElasticHook) Cancel() {
	hook.ctxCancel()
}
