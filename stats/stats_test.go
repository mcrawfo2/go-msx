package stats

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
	"time"
)

func TestNewPushConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *PushConfig
		wantErr bool
	}{
		{
			name: "Default",
			args: args{
				cfg: configtest.NewStaticConfig(nil),
			},
			want: &PushConfig{
				Enabled:   false,
				Url:       "",
				JobName:   "go_msx",
				Frequency: 15 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"stats.push.enabled":   "true",
					"stats.push.url":       "http://zipkin:16161",
					"stats.push.job-name":  "charlie",
					"stats.push.frequency": "5s",
				}),
			},
			want: &PushConfig{
				Enabled:   true,
				Url:       "http://zipkin:16161",
				JobName:   "charlie",
				Frequency: 5 * time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPushConfigFromConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPushConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPushConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
