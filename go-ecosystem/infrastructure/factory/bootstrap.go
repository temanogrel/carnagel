package factory

import (
	"os"

	"strings"

	"fmt"
	"time"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/consul"
	"git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/logging"
	consul_api "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/recursionpharma/logrus-stack"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
)

func Bootstrap(application string) (ecosystem.ConsulClient, logrus.FieldLogger, string, error) {

	api, err := consul_api.NewClient(consul_api.DefaultConfig())
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "Failed to create consul api client")
	}

	hostname, err := getHostname()
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "Failed to retrieve a hostname for the host")
	}

	logger, err := getLogger(application, hostname, api)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "Failed to create a logging client")
	}

	return consul.NewClient(api, logger), logger, hostname, nil
}

func getHostname() (string, error) {

	if os.Getenv("HOSTNAME") != "" {
		return os.Getenv("HOSTNAME"), nil
	}

	h, err := os.Hostname()
	if err != nil {
		return "", err
	}

	if strings.Contains(h, "vee.bz") {
		return h, nil
	}

	// os.Hostname does not return the fqdn
	return h + ".vee.bz", nil
}

func getLogger(application, hostname string, consul *consul_api.Client) (logrus.FieldLogger, error) {

	baseLogger := &logrus.Logger{
		Out:       os.Stderr,
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{},
	}

	// log stack traces for panic, fatal and error logs
	baseLogger.Hooks.Add(
		logrus_stack.LogrusStackHook{
			CallerLevels: []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel},
			StackLevels:  []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel},
		},
	)

	logger := baseLogger.WithField("application", application)

	services, _, err := consul.Catalog().Service("elasticsearch", "", &consul_api.QueryOptions{})
	if err != nil {
		logger.WithError(err).Fatal("Failed to retrieve elasticsearch")
	}

	if len(services) == 0 {
		logger.Warn("missing elasticsearch in consul")

		return logger, nil
	}

	var addresses []string
	for _, service := range services {
		addresses = append(addresses, fmt.Sprintf("http://%s:%d", service.ServiceAddress, service.ServicePort))
	}

	client, err := elastic.NewClient(elastic.SetURL(addresses...))
	if err != nil {
		logger.WithError(err).Fatal("Failed to construct elasticsearch client")
	}

	hook, err := logging.NewElasticHook(client, hostname, logrus.DebugLevel, func() string {
		return fmt.Sprintf("%s-%s", application, time.Now().Format("2006-01-02"))
	})

	if err != nil {
		logger.WithError(err).Fatal("Failed to create elasticsearch hook for logger")
	}

	baseLogger.Hooks.Add(hook)

	return logger, nil
}
