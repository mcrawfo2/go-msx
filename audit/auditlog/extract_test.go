package auditlog

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"reflect"
	"testing"
)

func TestExtractRequestDetails(t *testing.T) {
	type args struct {
		req  *restful.Request
		host string
		port int
	}
	tests := []struct {
		name string
		args args
		want *RequestDetails
	}{
		{
			name: "Proxied",
			args: args{
				req:  &restful.Request{
					Request: &http.Request{
						RemoteAddr: "10.10.10.10",
						Proto: "http",
						Header: map[string][]string{
							XForwardedForHeader: {"192.168.2.1"},
						},
					},
				},
				host: "10.10.10.12",
				port: 8080,
			},
			want: &RequestDetails{
				Source:   "192.168.2.1",
				Protocol: "http",
				Host:     "10.10.10.12",
				Port:     "8080",
			},
		},
		{
			name: "Unproxied",
			args: args{
				req:  &restful.Request{
					Request: &http.Request{
						RemoteAddr: "192.168.2.1",
						Proto: "http",
						Header: map[string][]string{},
					},
				},
				host: "10.10.10.12",
				port: 8080,
			},
			want: &RequestDetails{
				Source:   "192.168.2.1",
				Protocol: "http",
				Host:     "10.10.10.12",
				Port:     "8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractRequestDetails(tt.args.req, tt.args.host, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractRequestDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}
