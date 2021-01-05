package consulprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewConsulLeaderElectionConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ConsulLeaderElectionConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"spring.application.name": "TestNewConsulLeaderElectionConfigFromConfig/Defaults",
				}),
			},
			want: &ConsulLeaderElectionConfig{
				Enabled:          false,
				DefaultMasterKey: "service/TestNewConsulLeaderElectionConfigFromConfig/Defaults/leader",
				LeaderProperties: []LeaderProperties{},
			},
		},
		{
			name: "Enabled",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"spring.application.name":        "TestNewConsulLeaderElectionConfigFromConfig/Enabled",
					"consul.leader.election.enabled": "true",
				}),
			},
			want: &ConsulLeaderElectionConfig{
				Enabled:          true,
				DefaultMasterKey: "service/TestNewConsulLeaderElectionConfigFromConfig/Enabled/leader",
				LeaderProperties: []LeaderProperties{},
			},
		},
		{
			name: "EmbeddedDefaults",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"spring.application.name":                         "TestNewConsulLeaderElectionConfigFromConfig/EmbeddedDefaults",
					"consul.leader.election.enabled":                  "true",
					"consul.leader.election.default-master-key":       "service/${spring.application.name}/leader",
					"consul.leader.election.leader-properties[0].key": "${consul.leader.election.defaultMasterKey}",
				}),
			},
			want: &ConsulLeaderElectionConfig{
				Enabled:          true,
				DefaultMasterKey: "service/TestNewConsulLeaderElectionConfigFromConfig/EmbeddedDefaults/leader",
				LeaderProperties: []LeaderProperties{
					{
						Key:             "service/TestNewConsulLeaderElectionConfigFromConfig/EmbeddedDefaults/leader",
						HeartBeatMillis: 2000,
						BusyWaitMillis:  5000,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConsulLeaderElectionConfigFromConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConsulLeaderElectionConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConsulLeaderElectionConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
