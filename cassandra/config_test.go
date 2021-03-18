package cassandra

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
	"time"
)

func TestNewClusterConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ClusterConfig
		wantErr bool
	}{
		{
			name: "StructDefaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{}),
			},
			want: &ClusterConfig{
				Enabled:            false,
				KeyspaceName:       "system",
				ContactPoints:      "localhost",
				Port:               9042,
				Username:           "cassandra",
				Password:           "cassandra",
				Timeout:            15 * time.Second,
				ConnectTimeout:     5 * time.Second,
				Consistency:        "LOCAL_QUORUM",
				FullConsistency:    "ONE",
				PersistentSessions: false,
				KeyspaceOptions: KeyspaceOptions{
					Replications:             []string{"datacenter1"},
					DefaultReplicationFactor: 1,
				},
			},
		},
		{
			name: "CustomOptions",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.data.cassandra.enabled":                                     "true",
					"spring.data.cassandra.keyspace-name":                               "default",
					"spring.data.cassandra.contact-points":                              "remote-vm",
					"spring.data.cassandra.port":                                        "9999",
					"spring.data.cassandra.username":                                    "msxuser",
					"spring.data.cassandra.password":                                    "password",
					"spring.data.cassandra.timeout":                                     "30s",
					"spring.data.cassandra.connect-timeout":                             "10s",
					"spring.data.cassandra.consistency":                                 "LOCAL_ONE",
					"spring.data.cassandra.full-consistency":                            "QUORUM",
					"spring.data.cassandra.persistent-sessions":                         "true",
					"spring.data.cassandra.keyspace-options.replications[0]":            "dc1:2",
					"spring.data.cassandra.keyspace-options.replications[1]":            "dc2:3",
					"spring.data.cassandra.keyspace-options.default-replication-factor": "3",
				}),
			},
			want: &ClusterConfig{
				Enabled:            true,
				KeyspaceName:       "default",
				ContactPoints:      "remote-vm",
				Port:               9999,
				Username:           "msxuser",
				Password:           "password",
				Timeout:            30 * time.Second,
				ConnectTimeout:     10 * time.Second,
				Consistency:        "LOCAL_ONE",
				FullConsistency:    "QUORUM",
				PersistentSessions: true,
				KeyspaceOptions: KeyspaceOptions{
					Replications:             []string{"dc1:2", "dc2:3"},
					DefaultReplicationFactor: 3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClusterConfigFromConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClusterConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClusterConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
