// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/kafka"
	"cto-github.cisco.com/NFV-BU/go-msx/redis"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func setupServerDocumentDeps(t *testing.T) {
	ctx := context.Background()
	ctx = configtest.ContextWithNewInMemoryConfig(ctx,
		map[string]string{
			"spring.cloud.stream.kafka.binder.pool.enabled": "false",
		})

	err := kafka.ConfigurePool(ctx)
	if err != nil && err != kafka.ErrDisabled {
		t.Error(err)
	}

	err = redis.ConfigurePool(ctx)
	if err != nil && err != redis.ErrDisabled {
		t.Error(err)
	}
}

func TestServerDocumentor_Document(t *testing.T) {
	setupServerDocumentDeps(t)

	type fields struct {
		skip    bool
		server  *Server
		mutator ServerMutator
	}
	type args struct {
		binder string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Skip",
			fields: fields{
				skip: true,
			},
			args: args{
				binder: "",
			},
			wantErr: false,
		},
		{
			name:    "NoServer",
			fields:  fields{},
			args:    args{},
			wantErr: false,
		},
		{
			name: "Mutator",
			fields: fields{
				mutator: func(server *Server) {
					server.Description = types.NewStringPtr("mutated")
				},
			},
			args:    args{},
			wantErr: false,
		},
		{
			name:    "Binder-Kafka",
			fields:  fields{},
			args:    args{binder: "kafka"},
			wantErr: false,
		},
		{
			name:    "Binder-Redis",
			fields:  fields{},
			args:    args{binder: "redis"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ServerDocumentor{
				skip:    tt.fields.skip,
				server:  tt.fields.server,
				mutator: tt.fields.mutator,
			}
			if err := d.Document(tt.args.binder); (err != nil) != tt.wantErr {
				t.Errorf("Document() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerDocumentor_WithServerItem(t *testing.T) {
	server := &Server{
		MapOfAnything: map[string]interface{}{
			"key": "value",
		},
	}

	want := &ServerDocumentor{
		server: server,
	}

	got := new(ServerDocumentor).WithServerItem(server)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestServerDocumentor_WithServerMutator(t *testing.T) {
	mutator := func(server *Server) {}

	want := &ServerDocumentor{
		mutator: mutator,
	}

	got := new(ServerDocumentor).WithServerMutator(mutator)
	assert.True(t,
		reflect.DeepEqual(
			fmt.Sprintf("%p", want.mutator),
			fmt.Sprintf("%p", got.mutator),
		),
		testhelpers.Diff(want, got))
}

func TestServerDocumentor_WithSkip(t *testing.T) {
	skip := true

	want := &ServerDocumentor{
		skip: skip,
	}

	got := new(ServerDocumentor).WithSkip(skip)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestServerDocumentor_documentKafka(t *testing.T) {
	setupServerDocumentDeps(t)

	tests := []struct {
		name    string
		want    *Server
		wantErr bool
	}{
		{
			name: "Success",
			want: &Server{
				URL:         types.NewStringPtr("localhost:9092"),
				Description: types.NewStringPtr("CPX Internal Kafka"),
				Protocol:    types.NewStringPtr("kafka-plaintext"),
				Security: []map[string][]string{
					{
						"cpx": []string{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{}

			d := &ServerDocumentor{
				server: server,
			}

			d.documentKafka(server)
			assert.True(t,
				reflect.DeepEqual(tt.want, server),
				testhelpers.Diff(tt.want, server))
		})
	}
}

func TestServerDocumentor_documentRedis(t *testing.T) {
	setupServerDocumentDeps(t)

	tests := []struct {
		name    string
		want    *Server
		wantErr bool
	}{
		{
			name: "Success",
			want: &Server{
				URL:         types.NewStringPtr("localhost:6379"),
				Description: types.NewStringPtr("CPX Internal Redis"),
				Protocol:    types.NewStringPtr("redis"),
				Security: []map[string][]string{
					{
						"cpx": []string{},
					},
				},
				Bindings: &BindingsObject{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{}

			d := &ServerDocumentor{
				server: server,
			}

			d.documentRedis(server)
			assert.True(t,
				reflect.DeepEqual(tt.want, server),
				testhelpers.Diff(tt.want, server))
		})
	}
}
