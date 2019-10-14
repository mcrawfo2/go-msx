package cassandra

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const (
	configRootCassandraCluster = "spring.data.cassandra"
)

var (
	ErrDisabled = errors.New("Cassandra connection disabled")
	logger      = log.NewLogger("msx.cassandra")
)

type ClusterConfig struct {
	Enabled           bool          `config:"default=true"`
	KeyspaceName      string        // No default
	ContactPoints     string        `config:"default=localhost"` // comma separated
	Port              int           `config:"default=8500"`
	Username          string        `config:"default=cassandra"`
	Password          string        `config:"default=cassandra"`
	Timeout           time.Duration `config:"default=15s"`
	Consistency       string        `config:"default=LOCAL_QUORUM"`
	DataCenter        string        `config:"default=datacenter1"`
	ReplicationFactor int           `config:"default=1"`
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

type Cluster struct {
	config  *ClusterConfig
	cluster *gocql.ClusterConfig
}

func (c *Cluster) CreateSession() (*gocql.Session, error) {
	return c.cluster.CreateSession()
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

	//Configure authentication options if credentials available
	if clusterConfig.Username != "" && clusterConfig.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		}
	}

	return &Cluster{
		config: clusterConfig,
		cluster: cluster,
	}, nil
}

func NewClusterFromConfig(cfg *config.Config) (*Cluster, error) {
	clusterConfig := &ClusterConfig{}
	if err := cfg.Populate(clusterConfig, configRootCassandraCluster); err != nil {
		return nil, err
	}

	return NewCluster(clusterConfig)
}
