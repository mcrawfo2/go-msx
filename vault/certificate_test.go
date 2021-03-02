package vault

import (
	"reflect"
	"testing"
	"time"
)

func TestIssueCertificateRequest_Data(t *testing.T) {
	type fields struct {
		CommonName string
		Ttl        time.Duration
		AltNames   []string
		IpSans     []string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "Simple",
			fields: fields{
				CommonName: "TestIssueCertificateRequest_Data",
				Ttl:        30 * time.Minute,
				AltNames: []string{
					"TestIssueCertificateRequest",
					"localhost",
				},
				IpSans: []string{
					"127.0.0.1",
				},
			},
			want: map[string]interface{}{
				"common_name": "TestIssueCertificateRequest_Data",
				"ttl":         "30m0s",
				"alt_names":   "TestIssueCertificateRequest,localhost",
				"ip_sans":     "127.0.0.1",
			},
		},
		{
			name: "EmptyLists",
			fields: fields{
				CommonName: "TestIssueCertificateRequest_Data",
				Ttl:        30 * time.Minute,
			},
			want: map[string]interface{}{
				"common_name": "TestIssueCertificateRequest_Data",
				"ttl":         "30m0s",
				"alt_names":   "",
				"ip_sans":     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := IssueCertificateRequest{
				CommonName: tt.fields.CommonName,
				Ttl:        tt.fields.Ttl,
				AltNames:   tt.fields.AltNames,
				IpSans:     tt.fields.IpSans,
			}
			if got := r.Data(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data() = %v, want %v", got, tt.want)
			}
		})
	}
}
