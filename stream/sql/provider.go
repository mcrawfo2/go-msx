package sql

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	gosql "database/sql"
	"github.com/ThreeDotsLabs/watermill-sql/pkg/sql"
	"github.com/pkg/errors"
	"sync"
)

const (
	providerNameSql = "sql"
)

var ErrDisabled = sqldb.ErrDisabled
var loggerWatermillSql = log.NewLogger("watermill.sql")
var loggerAdapter = stream.NewWatermillLoggerAdapter(loggerWatermillSql)

var db *gosql.DB
var dbMtx sync.Mutex

type Provider struct{}

func getDatabase() (result *gosql.DB, err error) {
	dbMtx.Lock()
	defer dbMtx.Unlock()

	if db == nil {
		db, err = sqldb.Pool().NewSqlConnection()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create SQL connection")
		}
	}

	return db, nil
}

func (p *Provider) newSqlPublisher() (*sql.Publisher, error) {
	sqlPublisherDatabase, err := getDatabase()
	if err != nil {
		return nil, err
	}

	sqlPublisherConfig := sql.PublisherConfig{
		SchemaAdapter:        sql.DefaultPostgreSQLSchema{},
		AutoInitializeSchema: true,
	}

	sqlPublisher, err := sql.NewPublisher(sqlPublisherDatabase, sqlPublisherConfig, loggerAdapter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create SQL publisher")
	}

	return sqlPublisher, nil
}

func (p *Provider) newSqlSubscriber(cfg *config.Config, name string, streamBinding *stream.BindingConfiguration) (*sql.Subscriber, error) {
	sqlSubscriberDatabase, err := getDatabase()
	if err != nil {
		return nil, err
	}

	bindingConfiguration, err := NewBindingConfigurationFromConfig(cfg, name, streamBinding)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create binding configuration")
	}

	sqlSubscriberConfig := sql.SubscriberConfig{
		ConsumerGroup:  bindingConfiguration.StreamBindingConfig.Group,
		PollInterval:   bindingConfiguration.Consumer.PollInterval,
		ResendInterval: bindingConfiguration.Consumer.ResendInterval,
		RetryInterval:  bindingConfiguration.Consumer.RetryInterval,
		BackoffManager: sql.NewDefaultBackoffManager(
			bindingConfiguration.Consumer.PollInterval,
			bindingConfiguration.Consumer.RetryInterval),
		SchemaAdapter:    sql.DefaultPostgreSQLSchema{},
		OffsetsAdapter:   sql.DefaultPostgreSQLOffsetsAdapter{},
		InitializeSchema: true,
	}

	sqlSubscriber, err := sql.NewSubscriber(sqlSubscriberDatabase, sqlSubscriberConfig, loggerAdapter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create SQL subscriber")
	}

	return sqlSubscriber, nil
}

func (p *Provider) NewPublisher(_ *config.Config, _ string, streamBinding *stream.BindingConfiguration) (stream.Publisher, error) {
	sqlPublisher, err := p.newSqlPublisher()
	if err != nil {
		return nil, err
	}

	publisher := stream.NewTopicPublisher(sqlPublisher, streamBinding)
	return publisher, nil
}

func (p *Provider) NewSubscriber(cfg *config.Config, name string, streamBinding *stream.BindingConfiguration) (stream.Subscriber, error) {
	sqlSubscriber, err := p.newSqlSubscriber(cfg, name, streamBinding)
	if err != nil {
		return nil, err
	}

	return sqlSubscriber, nil
}

func RegisterProvider(cfg *config.Config) error {
	sqlConfig, err := sqldb.NewSqlConfigFromConfig(cfg)
	if err != nil {
		return err
	}

	if !sqlConfig.Enabled {
		return ErrDisabled
	}

	stream.RegisterProvider(providerNameSql, &Provider{})
	return nil
}
