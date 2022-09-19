package streamops_test

import (
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	_ "cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMessageValidator_ValidateMessage(t *testing.T) {
	type bodyStruct struct {
		A string `json:"a"`
	}

	portReflector := streamops.PortReflector{}

	tests := []struct {
		name       string
		portStruct interface{}
		payload    string
		metadata   map[string]string
		wantErr    bool
	}{
		{
			name: "Simple",
			portStruct: struct {
				Header string     `in:"header=Content-Type"`
				Body   bodyStruct `in:"body"`
			}{},
			metadata: map[string]string{
				"Content-Type": "application/json",
			},
			payload: `{"a":"value"}`,
		},
		{
			name: "ParseError",
			portStruct: struct {
				Header string     `in:"header=Content-Type"`
				Body   bodyStruct `in:"body"`
			}{},
			metadata: map[string]string{
				"Content-Type": "application/json",
			},
			payload: `{"a":"value"`,
			wantErr: true,
		},
		{
			name: "NumberValue",
			portStruct: struct {
				Header float64    `in:"header=X-Pi"`
				Body   bodyStruct `in:"body"`
			}{},
			metadata: map[string]string{
				"X-Pi": "3.14",
			},
			payload: `{"a":"value"}`,
		},
		{
			name: "IntegerValue",
			portStruct: struct {
				Timeout int        `in:"header=X-Backoff"`
				Body    bodyStruct `in:"body"`
			}{},
			metadata: map[string]string{
				"X-Backoff": "30",
			},
			payload: `{"a":"value"}`,
		},
		{
			name: "BooleanValue",
			portStruct: struct {
				Canary bool       `in:"header=X-Canary"`
				Body   bodyStruct `in:"body"`
			}{},
			metadata: map[string]string{
				"X-Canary": "false",
			},
			payload: `{"a":"value"}`,
		},
		{
			name: "ValidationError",
			portStruct: struct {
				Header string     `in:"header=Content-Type"`
				Body   bodyStruct `in:"body"`
			}{},
			metadata: map[string]string{},
			payload:  `{"a":"value"}`,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streamops.RegisterPortFieldValidationSchemaFunc(asyncapi.GetJsonValidationSchema)
			port, err := portReflector.ReflectInputPort(reflect.TypeOf(tt.portStruct))
			assert.NoError(t, err)

			msg := message.NewMessage(types.MustNewUUID().String(), []byte(tt.payload))
			msg.Metadata = message.Metadata(tt.metadata)
			msg.Payload = message.Payload(tt.payload)
			source := streamops.NewMessageDataSource("MY_TOPIC", msg)
			decoder := streamops.NewMessageDecoder(source, httpclient.MimeTypeApplicationJson, "")
			gotErr := streamops.NewMessageValidator(port, decoder).ValidateMessage()
			assert.Equal(t, tt.wantErr, gotErr != nil)
		})
	}

}
