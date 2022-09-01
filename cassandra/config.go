// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cassandra

import (
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const configRootCassandraCluster = "spring.data.cassandra"

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
	Disconnected       bool          `config:"default=${cli.flag.disconnected:false}"`
	KeyspaceName       string        `config:"default=system"`
	ContactPoints      string        `config:"default=localhost"` // comma separated
	Port               int           `config:"default=9042"`
	Username           string        `config:"default=cassandra"`
	Password           string        `config:"default=cassandra"`
	Timeout            time.Duration `config:"default=15s"`
	ConnectTimeout     time.Duration `config:"default=5s"`
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
