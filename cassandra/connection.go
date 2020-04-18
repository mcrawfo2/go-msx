package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

const (
	configRootCassandraCluster = "spring.data.cassandra"
	keyspaceSystem             = "system"
)

var (
	ErrDisabled = errors.New("Cassandra connection disabled")
	ErrNotFound = gocql.ErrNotFound
	logger      = log.NewLogger("msx.cassandra")
	gocqlLogger = log.NewLogger("gocql").Level(log.ErrorLevel)
)

func init() {
	gocql.Logger = gocqlLogger
}

type KeyspaceOptions struct {
	Replications             []string `config:"default=datacenter1"`
	DefaultReplicationFactor int      `config:"default=1"`
}

func (o KeyspaceOptions) ReplicationOptions() map[string]string {
	result := make(map[string]string)
	for _, replication := range o.Replications {
		parts := strings.SplitN(replication, ":", 2)
		if len(parts) == 1 || parts[1] == "" {
			result[parts[0]] = strconv.Itoa(o.DefaultReplicationFactor)
		} else {
			result[parts[0]] = parts[1]
		}
	}

	result[ddl.ReplicationOptionsKeyClass] = ddl.ClassNetworkTopologyStrategy
	return result
}

type ClusterConfig struct {
	Enabled            bool          `config:"default=false"`
	KeyspaceName       string        `config:"default=system"`
	ContactPoints      string        `config:"default=localhost"` // comma separated
	Port               int           `config:"default=9042"`
	Username           string        `config:"default=cassandra"`
	Password           string        `config:"default=cassandra"`
	Timeout            time.Duration `config:"default=15s"`
	Consistency        string        `config:"default=LOCAL_QUORUM"`
	FullConsistency    string        `config:"default=ONE"`
	PersistentSessions bool          `config:"default=false"`
	KeyspaceOptions    KeyspaceOptions
}

func (c ClusterConfig) Hosts() []string {
	hosts := strings.Split(c.ContactPoints, ",")
	for i, h := range hosts {
		hostParts := strings.SplitN(h, ":", 2)
		if len(hostParts) == 1 {
			hosts[i] = fmt.Sprintf("%s:%d", h, c.Port)
		}
	}
	return hosts
}

func NewClusterConfigFromConfig(cfg *config.Config) (*ClusterConfig, error) {
	clusterConfig := &ClusterConfig{}
	if err := cfg.Populate(clusterConfig, configRootCassandraCluster); err != nil {
		return nil, err
	}

	return clusterConfig, nil
}

type Cluster struct {
	config  *ClusterConfig
	cluster *gocql.ClusterConfig
}

func (c *Cluster) CreateSession() (*gocql.Session, error) {
	return c.cluster.CreateSession()
}

func (c *Cluster) FullConsistency() gocql.Consistency {
	return gocql.ParseConsistency(c.config.FullConsistency)
}

func (c *Cluster) createKeyspace(ctx context.Context, name string, options KeyspaceOptions) error {
	keyspaceQueryBuilder := new(ddl.KeyspaceQueryBuilder)
	session, err := c.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	keyspace := ddl.Keyspace{
		Name:               name,
		ReplicationOptions: options.ReplicationOptions(),
		DurableWrites:      true,
	}

	return session.
		Query(keyspaceQueryBuilder.CreateKeyspace(keyspace, true)).
		Consistency(c.FullConsistency()).
		WithContext(ctx).
		Exec()
}

func NewCluster(clusterConfig *ClusterConfig) (*Cluster, error) {
	if !clusterConfig.Enabled {
		logger.Warn("Cassandra connection disabled")
		return nil, ErrDisabled
	}

	cluster := gocql.NewCluster(clusterConfig.Hosts()...)
	cluster.Timeout = clusterConfig.Timeout
	cluster.Keyspace = clusterConfig.KeyspaceName
	cluster.Consistency = gocql.ParseConsistency(clusterConfig.Consistency)

	statsObserver := &StatsObserver{}
	traceObserver := &TraceObserver{}
	compositeQueryObserver := NewCompositeQueryObserver(statsObserver, traceObserver)
	compositeBatchObserver := NewCompositeBatchObserver(statsObserver, traceObserver)

	cluster.ConnectObserver = statsObserver
	cluster.QueryObserver = compositeQueryObserver
	cluster.BatchObserver = compositeBatchObserver

	//Configure authentication options if credentials available
	if clusterConfig.Username != "" && clusterConfig.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		}
	}

	return &Cluster{
		config:  clusterConfig,
		cluster: cluster,
	}, nil
}

func NewClusterFromConfig(cfg *config.Config) (*Cluster, error) {
	clusterConfig, err := NewClusterConfigFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return NewCluster(clusterConfig)
}

func NewSystemClusterFromConfig(cfg *config.Config) (*Cluster, error) {
	clusterConfig, err := NewClusterConfigFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	clusterConfig.KeyspaceName = keyspaceSystem

	return NewCluster(clusterConfig)
}
