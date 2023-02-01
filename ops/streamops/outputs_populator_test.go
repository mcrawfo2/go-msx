// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestOutputsPopulator_PopulateOutputs(t *testing.T) {
	type response struct {
		C string `json:"c"`
	}

	type outputs struct {
		A string                 `out:"header"`
		B response               `out:"body"`
		M types.Optional[string] `out:"header=contentType,optional"`
		N *string                `out:"header=contentEncoding,optional"`
		O types.Optional[string] `out:"channel,default=MY_TOPIC"`
		P types.Optional[string] `out:"messageId"`
	}

	port, err := PortReflector{}.ReflectOutputPort(reflect.TypeOf(outputs{}))
	assert.NoError(t, err)

	tests := []struct {
		name            string
		contentType     string
		contentEncoding string
		outputs         outputs
		wantErr         bool
		wantMetadata    map[string]string
		wantPayload     string
		wantChannel     string
		wantUuid        string
	}{
		{
			name:        "SimpleJson",
			contentType: "application/json",
			outputs: outputs{
				A: "value-a",
				B: response{
					C: "value-c",
				},
			},
			wantErr: false,
			wantMetadata: map[string]string{
				"a":           "value-a",
				"contentType": "application/json",
			},
			wantPayload: `{"c":"value-c"}` + "\n",
			wantChannel: "MY_TOPIC",
		},
		{
			name:        "OverrideJson",
			contentType: "text/plain",
			outputs: outputs{
				A: "value-a",
				B: response{
					C: "value-c",
				},
				M: types.OptionalOf("application/json"),
			},
			wantErr: false,
			wantMetadata: map[string]string{
				"a":           "value-a",
				"contentType": "application/json",
			},
			wantPayload: `{"c":"value-c"}` + "\n",
			wantChannel: "MY_TOPIC",
		},
		{
			name:        "OverrideChannel",
			contentType: "application/json",
			outputs: outputs{
				A: "value-a",
				B: response{
					C: "value-c",
				},
				O: types.OptionalOf("OTHER_TOPIC"),
			},
			wantErr: false,
			wantMetadata: map[string]string{
				"a":           "value-a",
				"contentType": "application/json",
			},
			wantPayload: `{"c":"value-c"}` + "\n",
			wantChannel: "OTHER_TOPIC",
		},
		{
			name:        "OverrideMessageId",
			contentType: "application/json",
			outputs: outputs{
				A: "value-a",
				B: response{
					C: "value-c",
				},
				P: types.OptionalOf("my-message-id"),
			},
			wantErr: false,
			wantMetadata: map[string]string{
				"a":           "value-a",
				"contentType": "application/json",
			},
			wantPayload: `{"c":"value-c"}` + "\n",
			wantChannel: "MY_TOPIC",
			wantUuid:    "my-message-id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var iface interface{} = tt.outputs
			var sink = new(MessageDataSink)
			p := OutputsPopulator{
				Outputs:         &iface,
				OutputPort:      port,
				ContentType:     tt.contentType,
				ContentEncoding: tt.contentEncoding,
				Encoder:         NewWatermillMessageEncoder(sink),
			}

			err = p.PopulateOutputs()
			assert.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			if !tt.wantErr {
				gotChannel, gotMessage := sink.Message(nil)
				assert.True(t,
					reflect.DeepEqual(tt.wantChannel, gotChannel.Value()),
					testhelpers.Diff(tt.wantChannel, gotChannel.Value()))
				assert.True(t,
					reflect.DeepEqual(tt.wantMetadata, map[string]string(gotMessage.Metadata)),
					testhelpers.Diff(tt.wantMetadata, map[string]string(gotMessage.Metadata)))
				assert.True(t,
					reflect.DeepEqual([]byte(tt.wantPayload), []byte(gotMessage.Payload)),
					testhelpers.Diff([]byte(tt.wantPayload), []byte(gotMessage.Payload)))
				if tt.wantUuid != "" {
					assert.True(t,
						reflect.DeepEqual(tt.wantUuid, gotMessage.UUID),
						testhelpers.Diff(tt.wantUuid, gotMessage.UUID))
				}
			}
		})
	}
}
