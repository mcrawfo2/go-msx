package stream

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterByMetaData(t *testing.T) {
	type args struct {
		key    string
		values []string
	}
	tests := []struct {
		name         string
		args         args
		wantPositive []message.Metadata
		wantNegative []message.Metadata
	}{
		{
			name:         "Single",
			args:         args{
				key:    "EntityType",
				values: []string{"TENANT"},
			},
			wantPositive: []message.Metadata{{"EntityType": "TENANT"}},
			wantNegative: []message.Metadata{{"EntityType": "SERVICE_INSTANCE"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := FilterByMetaData(tt.args.key, tt.args.values...)
			for _, meta := range tt.wantPositive {
				assert.True(t, filter(context.Background(), meta))
			}
			for _, meta := range tt.wantNegative {
				assert.False(t, filter(context.Background(), meta))
			}
		})
	}
}

func TestFilterMessage(t *testing.T) {
	type args struct {
		msg     *message.Message
		filters []MessageFilter
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "SingleFilterTrue",
			args: args{
				msg:     &message.Message{
					Metadata: message.Metadata{
						"EntityType": "TENANT",
					},
				},
				filters: []MessageFilter{
					FilterByMetaData("EntityType", "TENANT"),
				},
			},
			want: true,
		},
		{
			name: "SingleFilterFalse",
			args: args{
				msg:     &message.Message{
					Metadata: message.Metadata{
						"EntityType": "TENANT",
					},
				},
				filters: []MessageFilter{
					FilterByMetaData("EntityType", "SERVICE_INSTANCE"),
				},
			},
			want: false,
		},
		{
			name: "MultiFilterTrue",
			args: args{
				msg:     &message.Message{
					Metadata: message.Metadata{
						"EntityType": "TENANT",
						"Status": "CREATED",
					},
				},
				filters: []MessageFilter{
					FilterByMetaData("EntityType", "TENANT"),
					FilterByMetaData("Status", "CREATED", "DELETED"),
				},
			},
			want: true,
		},
		{
			name: "MultiFilterFalse",
			args: args{
				msg:     &message.Message{
					Metadata: message.Metadata{
						"EntityType": "TENANT",
						"Status": "CREATED",
					},
				},
				filters: []MessageFilter{
					FilterByMetaData("EntityType", "TENANT"),
					FilterByMetaData("Status", "DELETED"),
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterMessage(tt.args.msg, tt.args.filters)
			assert.Equal(t, tt.want, got)
		})
	}
}
