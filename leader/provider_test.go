package leader

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRegisterLeadershipProvider(t *testing.T) {
	type args struct {
		provider LeadershipProvider
	}
	tests := []struct {
		name string
		args args
		wantRegistered bool
	}{
		{
			name: "Nil",
			args: args{},
			wantRegistered: false,
		},
		{
			name: "NotNil",
			args: args{
				provider: new(MockLeadershipProvider),
			},
			wantRegistered: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leadershipProvider = nil
			RegisterLeadershipProvider(tt.args.provider)
			assert.True(t, tt.wantRegistered == IsLeadershipProviderRegistered())
		})
	}
}

func TestIsLeader(t *testing.T) {
	provider := new(MockLeadershipProvider)
	provider.
		On("IsLeader", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).
		Return(true).
		Times(1)

	type args struct {
		provider LeadershipProvider
	}
	tests := []struct {
		name string
		args args
		want bool
		wantErr bool
	}{
		{
			name: "Nil",
			args: args{},
			want: false,
			wantErr: true,
		},
		{
			name: "NotNil",
			args: args{
				provider: provider,
			},
			want: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leadershipProvider = nil
			RegisterLeadershipProvider(tt.args.provider)

			got, gotErr := IsLeader(context.Background(), "some-key")
			assert.Equal(t, tt.want, got)
			assert.True(t, tt.wantErr == (gotErr != nil))
		})
	}

	mock.AssertExpectationsForObjects(t, provider)
}

func TestIsMasterLeader(t *testing.T) {
	provider := new(MockLeadershipProvider)
	provider.
		On("IsLeader", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).
		Return(true).
		Times(1)
	provider.
		On("MasterKey", mock.AnythingOfType("*context.emptyCtx")).
		Return("some-key").
		Times(1)

	type args struct {
		provider LeadershipProvider
	}
	tests := []struct {
		name string
		args args
		want bool
		wantErr bool
	}{
		{
			name: "Nil",
			args: args{},
			want: false,
			wantErr: true,
		},
		{
			name: "NotNil",
			args: args{
				provider: provider,
			},
			want: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leadershipProvider = nil
			RegisterLeadershipProvider(tt.args.provider)

			got, gotErr := IsMasterLeader(context.Background())
			assert.Equal(t, tt.want, got)
			assert.True(t, tt.wantErr == (gotErr != nil))
		})
	}

	mock.AssertExpectationsForObjects(t, provider)
}

func TestReleaseLeadership(t *testing.T) {
	provider := new(MockLeadershipProvider)
	provider.
		On("ReleaseLeadership", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).
		Times(1)

	type args struct {
		provider LeadershipProvider
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{
			name: "Nil",
			args: args{},
			wantErr: true,
		},
		{
			name: "NotNil",
			args: args{
				provider: provider,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leadershipProvider = nil
			RegisterLeadershipProvider(tt.args.provider)

			gotErr := ReleaseLeadership(context.Background(), "some-key")
			assert.True(t, tt.wantErr == (gotErr != nil))
		})
	}

	mock.AssertExpectationsForObjects(t, provider)
}

func TestReleaseMasterLeadership(t *testing.T) {
	provider := new(MockLeadershipProvider)
	provider.
		On("ReleaseLeadership", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).
		Times(1)
	provider.
		On("MasterKey", mock.AnythingOfType("*context.emptyCtx")).
		Return("some-key").
		Times(1)

	type args struct {
		provider LeadershipProvider
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{
			name: "Nil",
			args: args{},
			wantErr: true,
		},
		{
			name: "NotNil",
			args: args{
				provider: provider,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leadershipProvider = nil
			RegisterLeadershipProvider(tt.args.provider)

			gotErr := ReleaseMasterLeadership(context.Background())
			assert.True(t, tt.wantErr == (gotErr != nil))
		})
	}

	mock.AssertExpectationsForObjects(t, provider)
}

func TestStart(t *testing.T) {
	provider := new(MockLeadershipProvider)
	provider.
		On("Start", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil).
		Times(1)

	type args struct {
		provider LeadershipProvider
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{
			name: "Nil",
			args: args{},
			wantErr: true,
		},
		{
			name: "NotNil",
			args: args{
				provider: provider,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leadershipProvider = nil
			RegisterLeadershipProvider(tt.args.provider)

			gotErr := Start(context.Background())
			assert.True(t, tt.wantErr == (gotErr != nil))
		})
	}

	mock.AssertExpectationsForObjects(t, provider)
}

func TestStop(t *testing.T) {
	provider := new(MockLeadershipProvider)
	provider.
		On("Stop", mock.AnythingOfType("*context.emptyCtx")).
		Return(nil).
		Times(1)

	type args struct {
		provider LeadershipProvider
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{
			name: "Nil",
			args: args{},
			wantErr: true,
		},
		{
			name: "NotNil",
			args: args{
				provider: provider,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leadershipProvider = nil
			RegisterLeadershipProvider(tt.args.provider)

			gotErr := Stop(context.Background())
			assert.True(t, tt.wantErr == (gotErr != nil))
		})
	}

	mock.AssertExpectationsForObjects(t, provider)
}
