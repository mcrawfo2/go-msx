package vault

import (
	"context"
	"crypto/tls"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func Test_newTraceConnection(t *testing.T) {
	mockConnection := new(MockConnection)

	type args struct {
		api ConnectionApi
	}
	tests := []struct {
		name string
		args args
		want traceConnection
	}{
		{
			name: "Success",
			args: args{
				api: mockConnection,
			},
			want: traceConnection{
				ConnectionApi: mockConnection,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTraceConnection(tt.args.api); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTraceConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_traceConnection_CreateTransitKey(t *testing.T) {
	type fields struct {
		ConnectionApi ConnectionApi
	}
	type args struct {
		ctx     context.Context
		keyName string
		request CreateTransitKeyRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("CreateTransitKey",
							mock.AnythingOfType("*context.valueCtx"),
							"my-key",
							mock.AnythingOfType("CreateTransitKeyRequest")).
						Return(nil).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:     context.Background(),
				keyName: "my-key",
				request: CreateTransitKeyRequest{},
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("CreateTransitKey",
							mock.AnythingOfType("*context.valueCtx"),
							"my-key",
							mock.AnythingOfType("CreateTransitKeyRequest")).
						Return(errors.New("error")).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:     context.Background(),
				keyName: "my-key",
				request: CreateTransitKeyRequest{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := traceConnection{
				ConnectionApi: tt.fields.ConnectionApi,
			}
			if err := s.CreateTransitKey(tt.args.ctx, tt.args.keyName, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("CreateTransitKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.fields.ConnectionApi.(*MockConnection).AssertExpectations(t)
		})
	}
}

func Test_traceConnection_Health(t *testing.T) {
	type fields struct {
		ConnectionApi ConnectionApi
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse *api.HealthResponse
		wantErr      bool
	}{
		{
			name: "Success",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("Health",
							mock.AnythingOfType("*context.valueCtx")).
						Return(&api.HealthResponse{
							ClusterName: "cluster-name",
						}, nil).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			wantResponse: &api.HealthResponse{
				ClusterName: "cluster-name",
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("Health",
							mock.AnythingOfType("*context.valueCtx")).
						Return(nil, errors.New("")).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			wantResponse: nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := traceConnection{
				ConnectionApi: tt.fields.ConnectionApi,
			}
			gotResponse, err := s.Health(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Health() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("Health() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func Test_traceConnection_IssueCertificate(t *testing.T) {
	type fields struct {
		ConnectionApi ConnectionApi
	}
	type args struct {
		ctx     context.Context
		role    string
		request IssueCertificateRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCert *tls.Certificate
		wantErr  bool
	}{
		{
			name: "Success",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("IssueCertificate",
							mock.AnythingOfType("*context.valueCtx"),
							"role",
							mock.AnythingOfType("IssueCertificateRequest")).
						Return(&tls.Certificate{}, nil).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:     context.Background(),
				role:    "role",
				request: IssueCertificateRequest{},
			},
			wantCert: &tls.Certificate{},
			wantErr:  false,
		},
		{
			name: "Error",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("IssueCertificate",
							mock.AnythingOfType("*context.valueCtx"),
							"role",
							mock.AnythingOfType("IssueCertificateRequest")).
						Return(nil, errors.New("")).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:     context.Background(),
				role:    "role",
				request: IssueCertificateRequest{},
			},
			wantCert: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := traceConnection{
				ConnectionApi: tt.fields.ConnectionApi,
			}
			gotCert, err := s.IssueCertificate(tt.args.ctx, tt.args.role, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("IssueCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCert, tt.wantCert) {
				t.Errorf("IssueCertificate() gotCert = %v, want %v", gotCert, tt.wantCert)
			}
		})
	}
}

func Test_traceConnection_ListSecrets(t *testing.T) {
	type fields struct {
		ConnectionApi ConnectionApi
	}
	type args struct {
		ctx  context.Context
		path string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantResults map[string]string
		wantErr     bool
	}{
		{
			name: "Success",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("ListSecrets",
							mock.AnythingOfType("*context.valueCtx"),
							"path").
						Return(map[string]string{}, nil).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:  context.Background(),
				path: "path",
			},
			wantResults: map[string]string{},
			wantErr:     false,
		},
		{
			name: "Error",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("ListSecrets",
							mock.AnythingOfType("*context.valueCtx"),
							"path").
						Return(nil, errors.New("")).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:  context.Background(),
				path: "path",
			},
			wantResults: nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := traceConnection{
				ConnectionApi: tt.fields.ConnectionApi,
			}
			gotResults, err := s.ListSecrets(tt.args.ctx, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListSecrets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("ListSecrets() gotResults = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func Test_traceConnection_Observe(t *testing.T) {
	t.Skipped()
}

func Test_traceConnection_RemoveSecrets(t *testing.T) {
	type fields struct {
		ConnectionApi ConnectionApi
	}
	type args struct {
		ctx  context.Context
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("RemoveSecrets",
							mock.AnythingOfType("*context.valueCtx"),
							"path").
						Return(nil).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:  context.Background(),
				path: "path",
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("RemoveSecrets",
							mock.AnythingOfType("*context.valueCtx"),
							"path").
						Return(errors.New("")).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:  context.Background(),
				path: "path",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := traceConnection{
				ConnectionApi: tt.fields.ConnectionApi,
			}
			if err := s.RemoveSecrets(tt.args.ctx, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("RemoveSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_traceConnection_StoreSecrets(t *testing.T) {
	type fields struct {
		ConnectionApi ConnectionApi
	}
	type args struct {
		ctx     context.Context
		path    string
		secrets map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("StoreSecrets",
							mock.AnythingOfType("*context.valueCtx"),
							"path",
							map[string]string{}).
						Return(nil).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:     context.Background(),
				path:    "path",
				secrets: map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("StoreSecrets",
							mock.AnythingOfType("*context.valueCtx"),
							"path",
							map[string]string{}).
						Return(errors.New("")).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:     context.Background(),
				path:    "path",
				secrets: map[string]string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := traceConnection{
				ConnectionApi: tt.fields.ConnectionApi,
			}
			if err := s.StoreSecrets(tt.args.ctx, tt.args.path, tt.args.secrets); (err != nil) != tt.wantErr {
				t.Errorf("StoreSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_traceConnection_TransitDecrypt(t *testing.T) {
	type fields struct {
		ConnectionApi ConnectionApi
	}
	type args struct {
		ctx        context.Context
		keyName    string
		ciphertext string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantPlaintext string
		wantErr       bool
	}{
		{
			name: "Success",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("TransitDecrypt",
							mock.AnythingOfType("*context.valueCtx"),
							"keyName",
							"ciphertext").
						Return("plaintext", nil).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:        context.Background(),
				keyName:    "keyName",
				ciphertext: "ciphertext",
			},
			wantPlaintext: "plaintext",
			wantErr:       false,
		},
		{
			name: "Error",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("TransitDecrypt",
							mock.AnythingOfType("*context.valueCtx"),
							"keyName",
							"ciphertext").
						Return("", errors.New("")).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:        context.Background(),
				keyName:    "keyName",
				ciphertext: "ciphertext",
			},
			wantPlaintext: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := traceConnection{
				ConnectionApi: tt.fields.ConnectionApi,
			}
			gotPlaintext, err := s.TransitDecrypt(tt.args.ctx, tt.args.keyName, tt.args.ciphertext)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransitDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPlaintext != tt.wantPlaintext {
				t.Errorf("TransitDecrypt() gotPlaintext = %v, want %v", gotPlaintext, tt.wantPlaintext)
			}
		})
	}
}

func Test_traceConnection_TransitEncrypt(t *testing.T) {
	type fields struct {
		ConnectionApi ConnectionApi
	}
	type args struct {
		ctx       context.Context
		keyName   string
		plaintext string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantCiphertext string
		wantErr        bool
	}{
		{
			name: "Success",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("TransitEncrypt",
							mock.AnythingOfType("*context.valueCtx"),
							"keyName",
							"plaintext").
						Return("ciphertext", nil).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:       context.Background(),
				keyName:   "keyName",
				plaintext: "plaintext",
			},
			wantCiphertext: "ciphertext",
			wantErr:        false,
		},
		{
			name: "Error",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("TransitEncrypt",
							mock.AnythingOfType("*context.valueCtx"),
							"keyName",
							"plaintext").
						Return("", errors.New("")).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:       context.Background(),
				keyName:   "keyName",
				plaintext: "plaintext",
			},
			wantCiphertext: "",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := traceConnection{
				ConnectionApi: tt.fields.ConnectionApi,
			}
			gotCiphertext, err := s.TransitEncrypt(tt.args.ctx, tt.args.keyName, tt.args.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransitEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCiphertext != tt.wantCiphertext {
				t.Errorf("TransitEncrypt() gotCiphertext = %v, want %v", gotCiphertext, tt.wantCiphertext)
			}
		})
	}
}

func Test_traceConnection_GenerateRandomBytes(t *testing.T) {
	type fields struct {
		ConnectionApi ConnectionApi
	}
	type args struct {
		ctx    context.Context
		length int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("GenerateRandomBytes",
							mock.AnythingOfType("*context.valueCtx"),
							3).
						Return([]byte{1, 2, 3}, nil).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:    context.Background(),
				length: 3,
			},
			want: []byte{1,2,3},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				ConnectionApi: func() ConnectionApi {
					mockConnection := new(MockConnection)
					mockConnection.
						On("GenerateRandomBytes",
							mock.AnythingOfType("*context.valueCtx"),
							128).
						Return(nil, errors.New("")).
						Once()
					return mockConnection
				}(),
			},
			args: args{
				ctx:    context.Background(),
				length: 128,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := traceConnection{
				ConnectionApi: tt.fields.ConnectionApi,
			}
			if got, err := s.GenerateRandomBytes(tt.args.ctx, tt.args.length); (err != nil) != tt.wantErr {
				t.Errorf("StoreSecrets() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StoreSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}
