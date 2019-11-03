package main

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"net/http"
)

const (
	appName = "someservice"
	kafkaTopicExample = "EXAMPLE_TOPIC"
)

var logger = log.NewLogger(appName)

func init() {
	app.OnEvent(app.EventStart, app.PhaseBefore, addWebService)
	app.OnEvent(app.EventStart, app.PhaseDuring, dumpConfiguration)
	app.OnEvent(app.EventStart, app.PhaseDuring, subscribeExampleTopic)
	app.OnEvent(app.EventReady, app.PhaseDuring, listGauges)
	app.OnEvent(app.EventReady, app.PhaseDuring, sendExampleTopicMessage)
}

func dumpConfiguration(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	if cfg == nil {
		return errors.New("Failed to obtain application config")
	}
	quiet, _ := cfg.BoolOr("cli.flag.quiet", false)
	if !quiet {
		logger.Info("Dumping application configuration")
		cfg.Each(func(name, value string) {
			logger.Infof("%s: %s", name, value)
		})
	}
	return nil
}

func subscribeExampleTopic(ctx context.Context) error {
	return stream.AddListener(kafkaTopicExample, func(msg *message.Message) error {
		logger.WithContext(msg.Context()).WithField("messageId", msg.UUID).Infof("received message payload: %s", string(msg.Payload))
		return errors.New("some error occurred")
	})
}

func addWebService(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	svc := server.NewService()

	svc.Route(svc.GET("/test").
		Operation("retrieveTest").
		To(webservice.HttpHandlerController(myTestEndpoint)).
		Do(webservice.StandardRoute). // Authenticated, Returns 200/400/401/403
		Produces(webservice.MIME_JSON).
		Filter(webservice.PermissionsFilter(security.PermissionIsApiAdmin)))

	return nil
}

func myTestEndpoint(writer http.ResponseWriter, req *http.Request) {
	userContext := security.UserContextFromContext(req.Context())
	writer.WriteHeader(200)
	if _, err := writer.Write([]byte(`{"user":"` + userContext.UserName + `"}`)); err != nil {
		logger.WithError(err).Error("Failed to write result body")
	}
}

func listGauges(ctx context.Context) error {
	cassandraPool := cassandra.PoolFromContext(ctx)
	if cassandraPool == nil {
		return errors.New("Cassandra connection pool not found")
	}

	return cassandraPool.WithSession(listGaugesFromSession)
}

func listGaugesFromSession(session *gocql.Session) error {
	var serviceType, deviceType, deviceSubType, beatType *string
	if err := session.Query(`SELECT servicetype, devicetype, devicesubtype, beattype FROM gauges LIMIT 1 ALLOW FILTERING`).
			Consistency(gocql.One).
			Scan(&serviceType, &deviceType, &deviceSubType, &beatType); err != nil && err != gocql.ErrNotFound {
		logger.Error(err)
	} else if err != gocql.ErrNotFound {
		logger.Infof("Found gauges: serviceType=%s deviceType=%s deviceSubType=%s beatType=%s",
			*serviceType, *deviceType, *deviceSubType, *beatType)
	}
	return nil
}

func migrate(ctx context.Context) error {
	logger.Info("Migrate activity here")
	return nil
}

func populate(ctx context.Context) error {
	logger.Info("Populate activity here")
	return errors.New("Population failed")
}

func sendExampleTopicMessage(ctx context.Context) error {
	if err := stream.Publish(ctx, kafkaTopicExample, []byte("Test Message"), nil); err != nil {
		logger.Error(err)
	}
	return nil
}

func main() {
	cli.RootCmd().PersistentFlags().Bool("quiet", false, "Be quiet")
	if _, err := app.AddCommand("migrate", "Migrate database schema", migrate); err != nil {
		cli.Fatal(err)
	}
	if _, err := app.AddCommand("populate", "Populate remote microservices", populate); err != nil {
		cli.Fatal(err)
	}
	app.Run(appName)
}
