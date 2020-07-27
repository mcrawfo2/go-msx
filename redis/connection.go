package redis

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const (
	configRootRedis = "spring.redis"
)

var (
	ErrDisabled = errors.New("Redis connection disabled")
	ErrNotFound = redis.Nil
	logger      = log.NewLogger("msx.redis")
)

// RedisConfig represents Redis config options
type ConnectionConfig struct {
	Enable      bool   `config:"default=false"`
	Host        string `config:"default=localhost"`
	Port        int    `config:"default=6379"`
	Password    string `config:"default="`
	DB          int    `config:"default=0"`
	Sentinel    SentinelConfig
	MaxRetries  int   `config:"default=2"`
	IdleTimeout int64 `config:"default=1"`
}

func (c ConnectionConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type SentinelConfig struct {
	Enable bool     `config:"default=false"`
	Master string   `config:"default=mymaster"`
	Nodes  []string `config:"default=localhost:26379"`
}

type Connection struct {
	config  *ConnectionConfig
	client  *redis.Client
	version string
}

func (c *Connection) Client(ctx context.Context) *redis.Client {
	return c.client.WithContext(ctx)
}

func (c *Connection) Version() string {
	return c.version
}

func newSentinelClient(cfg *ConnectionConfig) *redis.Client {
	logger.Infof("Connecting to redis sentinel address: %v", cfg.Sentinel.Nodes)

	options := redis.FailoverOptions{
		Password:      cfg.Password,
		DB:            cfg.DB,
		MasterName:    cfg.Sentinel.Master,
		SentinelAddrs: cfg.Sentinel.Nodes,
		MaxRetries:    cfg.MaxRetries,
		IdleTimeout:   time.Duration(cfg.IdleTimeout) * time.Minute,
	}

	return redis.NewFailoverClient(&options)
}

func newStandaloneClient(cfg *ConnectionConfig) *redis.Client {
	address := cfg.Address()
	logger.Infof("Connecting to redis standalone address: %s", address)

	options := redis.Options{
		Addr:        address,
		Password:    cfg.Password,
		DB:          cfg.DB,
		MaxRetries:  cfg.MaxRetries,
		IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Minute,
	}

	return redis.NewClient(&options)
}

func NewConnection(cfg *ConnectionConfig) (*Connection, error) {
	var client *redis.Client

	if !cfg.Enable {
		return nil, ErrDisabled
	}

	if cfg.Sentinel.Enable {
		client = newSentinelClient(cfg)
	} else {
		client = newStandaloneClient(cfg)
	}

	if res, err := client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "Redis ping returned error")
	} else {
		logger.Info("Redis ping returned: ", res)
	}

	var version string
	if text, err := client.Info("server").Result(); err != nil {
		return nil, errors.Wrap(err, "Redis server info returned error")
	} else {
		var info = make(map[string]string)
		for _, line := range strings.Split(text, "\r\n") {
			if strings.HasPrefix(line, "#") || len(line) == 0 {
				continue
			}

			lineParts := strings.SplitN(line, ":", 2)
			if len(lineParts) == 2 {
				info[lineParts[0]] = lineParts[1]
			}
		}

		var ok bool
		if version, ok = info[`redis_version`]; ok {
			logger.Info("Redis server version: ", version)
		} else {
			version = "Unknown"
			logger.Warn("Redis server version unknown")
		}
	}

	client.AddHook(new(statsHook))
	client.AddHook(new(traceHook))

	return &Connection{
		config:  cfg,
		client:  client,
		version: version,
	}, nil
}

func NewConnectionFromConfig(cfg *config.Config) (*Connection, error) {
	connectionConfig := &ConnectionConfig{}
	if err := cfg.Populate(connectionConfig, configRootRedis); err != nil {
		return nil, err
	}

	return NewConnection(connectionConfig)
}
