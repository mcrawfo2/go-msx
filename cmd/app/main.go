package main

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/migrate"
	_ "cto-github.cisco.com/NFV-BU/go-msx/cassandra/migrate"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"net/http"
)

const (
	appName                = "someservice"
	kafkaTopicExample      = "EXAMPLE_TOPIC"
	kafkaTopicNotification = "NOTIFICATION_TOPIC"
	channelTopicLoopback   = "LOOPBACK_TOPIC"
	configKeyQuiet         = "cli.flag.quiet"
)

var logger = log.NewLogger(appName)

func init() {
	app.OnEvent(app.EventCommand, app.CommandRoot, func(ctx context.Context) error {
		app.OnEvent(app.EventStart, app.PhaseBefore, addWebService)
		app.OnEvent(app.EventStart, app.PhaseDuring, dumpConfiguration)
		app.OnEvent(app.EventStart, app.PhaseDuring, subscribeExampleTopic)
		app.OnEvent(app.EventStart, app.PhaseDuring, subscribeLoopbackTopic)
		app.OnEvent(app.EventReady, app.PhaseDuring, listGauges)
		app.OnEvent(app.EventReady, app.PhaseDuring, sendLoopbackTopicMessage)
		return nil
	})

	app.OnEvent(app.EventCommand, app.CommandMigrate, func(ctx context.Context) error {
		app.OnEvent(app.EventStart, app.PhaseDuring, addMigrations)
		return nil
	})
}

func addMigrations(ctx context.Context) error {
	manifest := migrate.ManifestFromContext(ctx)

	return types.ErrorList{
		manifest.AddCqlStringMigration("3.8.0.1", "Create first table", "CREATE TABLE first (value text PRIMARY KEY)"),
		manifest.AddCqlStringMigration("3.8.0.2", "Create second table", "CREATE TABLE second (value text PRIMARY KEY)"),
		manifest.AddCqlStringMigration("3.8.0.3", "Drop first table", "DROP TABLE first"),
		manifest.AddCqlFileMigration("3.8.0.4", "Create third table", "3.8.0/V3_8_0_4__CREATE_THIRD_TABLE.cql"),
	}.Filter()
}

func dumpConfiguration(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	if cfg == nil {
		return errors.New("Failed to obtain application config")
	}
	quiet, _ := cfg.BoolOr(configKeyQuiet, false)
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
		return errors.New("some error occurred")
	})
}

func subscribeLoopbackTopic(ctx context.Context) error {
	return stream.AddListener(channelTopicLoopback, func(msg *message.Message) error {
		logger.WithContext(msg.Context()).WithField("messageId", msg.UUID).Infof("received message payload: %s", string(msg.Payload))
		return nil
	})
}

func addWebService(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	svc, err := server.NewService("/api/v1/tests")
	if err != nil {
		return err
	}

	svc.Route(svc.GET("").
		Operation("queryTests").
		To(webservice.HttpHandlerController(myTestEndpoint)).
		Do(webservice.StandardReturns). // Returns 200/400/401/403
		Produces(webservice.MIME_JSON).
		Filter(webservice.PermissionsFilter(rbac.PermissionIsApiAdmin)))

	tenantIdParameter := svc.PathParameter("tenantId", "Tenant Id")

	svc.Route(svc.GET("/{tenantId}").
		Operation("queryTestsWithTenant").
		To(webservice.HttpHandlerController(myTestEndpoint)).
		Do(webservice.StandardReturns). // Returns 200/400/401/403
		Produces(webservice.MIME_JSON).
		Param(tenantIdParameter).
		Filter(webservice.PermissionsFilter(rbac.PermissionIsApiAdmin)).
		Filter(webservice.TenantFilter(tenantIdParameter)))

	return nil
}

func myTestEndpoint(writer http.ResponseWriter, req *http.Request) {
	sendNotificationTopicMessage(req.Context())
	sendLoopbackTopicMessage(req.Context())

	userContext := security.UserContextFromContext(req.Context())

	writer.WriteHeader(200)
	if _, err := writer.Write([]byte(`{"user":"` + userContext.UserName + `"}`)); err != nil {
		logger.WithError(err).Error("Failed to write result body")
	}
}

func listGauges(ctx context.Context) error {
	cassandraPool, err := cassandra.PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return cassandraPool.WithSession(func(session *gocql.Session) error {
		var serviceType, deviceType, deviceSubType, beatType *string
		if err := session.Query(`SELECT servicetype, devicetype, devicesubtype, beattype FROM msx_alerts.gauges LIMIT 1 ALLOW FILTERING`).
			WithContext(ctx).
			Consistency(gocql.One).
			Scan(&serviceType, &deviceType, &deviceSubType, &beatType); err != nil && err != gocql.ErrNotFound {
			logger.Error(err)
		} else if err != gocql.ErrNotFound {
			logger.Infof("Found gauges: serviceType=%s deviceType=%s deviceSubType=%s beatType=%s",
				*serviceType, *deviceType, *deviceSubType, *beatType)
		}
		return nil
	})
}

func populate(ctx context.Context, _ []string) error {
	logger.WithContext(ctx).Info("Populate activity here")
	return errors.New("Population failed")
}

func sendLoopbackTopicMessage(ctx context.Context) error {
	logger.WithContext(ctx).Infof("Sending test message to %q", channelTopicLoopback)
	if err := stream.Publish(ctx, channelTopicLoopback, []byte("Test Message"), nil); err != nil {
		return err
	}
	return nil
}

func sendNotificationTopicMessage(ctx context.Context) error {
	msg := []byte(`{"tenant":{"tenantId":"5a05f81d-3d9e-40b2-94b9-a2d2421f0de3","providerName":"CiscoSystems","tenantName":"5a05f81d-3d9e-40b2-94b9-a2d2421f0de3","tenantDescription":null,"vpnDescriptor":null,"displayName":"5a05f81d-3d9e-40b2-94b9-a2d2421f0de3","email":null,"phoneNumber":null,"url":"","tenantGroupName":null,"mobileNumber":null,"tenantExtension":null,"suspended":false},"provider":{"id":"fe3ad89c-449f-42f2-b4f8-b10ab7bc0266","email":"noreply@cisco.com","name":"CiscoSystems","notificationType":null,"locale":"en_US"},"restDetails":null,"statusMessage":null}`)
	if err := stream.Publish(ctx, kafkaTopicNotification, msg, nil); err != nil {
		return err
	}
	return nil
}

func main() {
	cli.RootCmd().PersistentFlags().Bool("quiet", false, "Be quiet")
	if _, err := app.AddCommand("populate", "Populate remote microservices", populate, app.Noop); err != nil {
		cli.Fatal(err)
	}
	app.Run(appName)
}
