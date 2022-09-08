// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/kafka"
	"cto-github.cisco.com/NFV-BU/go-msx/redis"
)

type ServerMutator func(server *Server)

type ServerDocumentor struct {
	skip    bool
	server  *Server
	mutator ServerMutator
}

func (d *ServerDocumentor) WithSkip(skip bool) *ServerDocumentor {
	d.skip = skip
	return d
}

func (d *ServerDocumentor) WithServerItem(server *Server) *ServerDocumentor {
	d.server = server
	return d
}

func (d *ServerDocumentor) WithServerMutator(fn ServerMutator) *ServerDocumentor {
	d.mutator = fn
	return d
}

func (d *ServerDocumentor) Document(binder string) (err error) {
	if d.skip {
		return nil
	}

	// Initialize
	server := d.server
	if server == nil {
		server = new(Server)
	}

	switch binder {
	case "kafka":
		d.documentKafka(server)
	case "redis":
		d.documentRedis(server)
	}

	// Mutator
	if d.mutator != nil {
		d.mutator(server)
	}

	// Publish
	Reflector.SpecEns().WithServersItem(binder, *server)

	return nil
}

func (d *ServerDocumentor) documentKafka(server *Server) {
	kafkaConfig := kafka.Pool().ConnectionConfig()

	server.WithURL(kafkaConfig.BrokerAddresses()[0])
	server.WithDescription("CPX Internal Kafka")

	if kafkaConfig.Tls.Enabled {
		server.WithProtocol("kafka-secure")
	} else {
		server.WithProtocol("kafka-plaintext")
	}

	server.WithSecurity(map[string][]string{
		"cpx": {},
	})
}

func (d *ServerDocumentor) documentRedis(server *Server) {
	redisConfig := redis.Pool().ConnectionConfig()

	server.WithURL(redisConfig.Address())
	server.WithDescription("CPX Internal Redis")
	server.WithProtocol("redis")
	server.WithBindings(*new(BindingsObject))
	server.WithSecurity(map[string][]string{
		"cpx": {},
	})
}
