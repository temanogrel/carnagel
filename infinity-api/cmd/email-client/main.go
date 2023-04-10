package main

import (
	"context"
	"flag"
	"fmt"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/infrastructure/factory"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	email_pkg "git.misc.vee.bz/carnagel/infinity-api/pkg/domain/email"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/payment"
	"git.misc.vee.bz/carnagel/infinity-api/pkg/domain/user"
	"github.com/go-pg/pg"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
	"os"
	"os/signal"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
	doneCh      = make(chan bool)
	hostname    string
	consul      ecosystem.ConsulClient
	logger      logrus.FieldLogger
	emailsSent  int

	app *infinity.Application
)

var (
	command = flag.String("command", "", "Command to run")
	email   = flag.String("email", "", "Specify email for the 'email-single-user' command")
	subject = flag.String("subject", "", "Subject of email")
	path    = flag.String("path", "", "Path to template file")
)

func init() {
	flag.Parse()

	var err error

	consul, logger, hostname, err = factory.Bootstrap("email-client")
	if err != nil {
		panic(err)
	}

	elasticClient, err := getElasticSearchClient()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to elasticsearch")
	}

	dbUser := consul.GetString("infinity/postgres/user", "infinity")
	dbName := consul.GetString("infinity/postgres/name", "infinity")
	dbPass := consul.GetString("infinity/postgres/pass", "infinity")
	dbHost := consul.GetString("infinity/postgres/host", "localhost")
	dbPort := consul.GetString("infinity/postgres/port", "5432")

	// use go.pg for the remainder of the time
	database := pg.Connect(&pg.Options{
		User:     dbUser,
		Addr:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		Password: dbPass,
		Database: dbName,
	})

	app = &infinity.Application{
		Logger:        logger,
		DB:            database,
		Consul:        consul,
		ElasticSearch: elasticClient,

		// Ecosystem clients
	}

	app.UserRepository = user.NewRepository(app)
	app.PaymentPlanRepository = payment.NewPaymentPlanRepository(app)
	app.EmailService = email_pkg.NewEmailService(
		app,
		consul.GetString("infinity/email/host", ""),
		int(consul.GetInt16("infinity/email/port", 465)),
		consul.GetString("infinity/email/username", ""),
		consul.GetString("infinity/email/password", ""),
		consul.GetString("infinity/email/fromAddress", ""),
	)
}

func main() {
	if *subject == "" {
		logger.Fatal("No subject provided")
	}

	if *path == "" {
		logger.Fatal("No path to template file provided")
	}

	criteria := infinity.NewUserRepositoryCriteria()
	criteria.Limit = 250
	criteria.Sorting = map[string]string{"created_at": "asc"}

	switch *command {
	case "email-premium":
		criteria.CurrentlyPremium = true

	case "email-ex-premiums":
		criteria.ExPremium = true

	case "email-never-premium":
		criteria.NeverPremium = true

	case "email-single-user":
		if *email == "" {
			logger.Fatal("No email specified")
		}

		app.EmailService.SendEmailFromTemplateFile(*email, *subject, *path)

		return

	case "email-all":
		// do nothing with criteria

	default:
		logger.WithField("command", *command).Fatal("Unknown command provided")
	}

	sendEmailToUsers(criteria)

	signalCh := make(chan os.Signal, 1)

	// Listen for shutdown
	signal.Notify(signalCh, os.Interrupt)

	select {
	case success := <-doneCh:
		app.Logger.
			WithFields(logrus.Fields{"success": success, "emailsSent": emailsSent}).
			Info("Finished running command")
		return

	case <-signalCh:
		app.Logger.Info("Received SIGINT, terminating...")

		// cleanup
		cancel()
	}
}

func sendEmailToUsers(criteria *infinity.UserRepositoryCriteria) {
	for {
		users, _, err := app.UserRepository.Matching(criteria)
		if err != nil {
			logger.WithError(err).Fatal("Failed to retrieve users by criteria")
		}

		for _, user := range users {
			app.EmailService.SendEmailFromTemplateFile(user.Email, *subject, *path)
		}

		emailsSent += len(users)

		if len(users) < criteria.Limit {
			doneCh <- true
			return
		}

		criteria.CreatedAfter = users[len(users)-1].CreatedAt
	}
}

func getElasticSearchClient() (*elastic.Client, error) {
	services, _, err := consul.API().Catalog().Service("elasticsearch", "", &api.QueryOptions{})
	if err != nil {
		logger.WithError(err).Fatal("Failed to retrieve elasticsearch")
	}

	var addresses []string
	for _, service := range services {
		addresses = append(addresses, fmt.Sprintf("http://%s:%d", service.ServiceAddress, service.ServicePort))
	}

	return elastic.NewClient(elastic.SetURL(addresses...))
}
