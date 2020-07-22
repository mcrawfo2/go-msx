package beater

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx-beats/meta"
	"cto-github.cisco.com/NFV-BU/go-msx-beats/state"
	"cto-github.cisco.com/NFV-BU/go-msx-beats/webconfig"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/templates/code/beats/api"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/elastic/beats/libbeat/common"
	"github.com/pkg/errors"
)

func init() {
	webconfig.SetDeviceDto(api.DeviceDTO{})
	webconfig.SetDeviceDtoList([]api.DeviceDTO{})
	webconfig.SetHostDto(Host{})
	webconfig.SetHostDtoList([]Host{})
}

type BeatState struct {
	Hosts []Host
}

type Host struct {
	Id        string   `config:"id" json:"id"`
	ServiceId string   `config:"serviceId" json:"serviceId"`
	Host      string   `config:"host" json:"host"`
	Tags      []string `config:"tags" json:"tags"`
}

func (h Host) NewEvent() common.MapStr {
	return map[string]interface{}{
		meta.FieldDeviceAddress: h.Host,
		meta.FieldDeviceId:      h.Id,
		meta.FieldServiceId:     h.ServiceId,
		meta.FieldTags:          h.Tags,
	}
}

func (h Host) NewLog() map[string]interface{} {
	return map[string]interface{}{
		trace.FieldDeviceAddress: h.Host,
		trace.FieldDeviceId:      h.Id,
		trace.FieldServiceId:     h.ServiceId,
	}
}

type BeatStateService struct {
	RunningState   *BeatState
	CandidateState *BeatState
	store          state.Store
	worker         *types.Worker
}

func (s *BeatStateService) Candidate(_ context.Context) interface{} {
	return s.CandidateState.Hosts
}

func (s *BeatStateService) Commit(ctx context.Context) error {
	logger.WithContext(ctx).Info("Committing candidate to running state")
	return s.worker.Run(func(_ context.Context) error {
		s.RunningState = &BeatState{
			Hosts: s.CandidateState.Hosts,
		}
		return s.store.Save(ctx, s.RunningState.Hosts)
	})
}

func (s *BeatStateService) Init(ctx context.Context) (err error) {
	logger.WithContext(ctx).Info("Loading running state")
	var data []byte
	var initState *BeatState

	if data, err = s.store.Get(ctx); err != nil {
		return errors.Wrap(err, "Failed to load state")
	}

	if initState, err = unmarshalState(data); err != nil {
		return errors.Wrap(err, "Failed to unmarshal state")
	}

	logger.WithContext(ctx).Debugf("Loaded state: %v", initState)

	s.RunningState = initState
	s.CandidateState = &BeatState{}
	return
}

func (s *BeatStateService) Running(ctx context.Context) interface{} {
	return s.RunningState.Hosts
}

func (s *BeatStateService) SetCandidate(ctx context.Context, data []byte) (err error) {
	logger.WithContext(ctx).Info("Setting candidate state")

	var candidateState *BeatState
	var devices []api.DeviceDTO

	if err = json.Unmarshal(data, &devices); err != nil {
		return errors.Wrap(err, "Failed to unmarshal devices")
	}

	if candidateState, err = mapDeviceStates(devices); err != nil {
		return errors.Wrap(err, "Failed to map device state")
	}

	s.CandidateState = candidateState
	return
}

func (s *BeatStateService) UnsetRunning(ctx context.Context) error {
	logger.WithContext(ctx).Info("Clearing running state")
	return s.worker.Run(func(_ context.Context) error {
		s.RunningState = &BeatState{Hosts: []Host{}}
		return s.store.Save(ctx, s.RunningState.Hosts)
	})
}

func (s *BeatStateService) ExtendRunning(ctx context.Context, data []byte) (err error) {
	logger.WithContext(ctx).Info("Extending running state to candidate state")

	var device api.DeviceDTO

	if err = json.Unmarshal(data, &device); err != nil {
		return errors.Wrap(err, "Failed to unmarshal device")
	}

	var candidateHost Host
	if candidateHost, err = mapDeviceState(device); err != nil {
		return errors.Wrap(err, "Failed to map device state")
	}

	return s.worker.Run(func(_ context.Context) error {
		var newState *BeatState
		if newState, err = s.copyRunningStateWithoutDevice(candidateHost.Id); err != nil {
			return errors.Wrap(err, "Failed to copy running state to candidate state")
		}

		newState.Hosts = append(newState.Hosts, candidateHost)
		s.RunningState = newState
		return s.store.Save(ctx, s.RunningState.Hosts)
	})
}

func (s *BeatStateService) ShrinkRunning(ctx context.Context, deviceId string) (err error) {
	logger.WithContext(ctx).Info("Shrinking running state to candidate state")

	return s.worker.Run(func(_ context.Context) error {
		var newState *BeatState
		if newState, err = s.copyRunningStateWithoutDevice(deviceId); err != nil {
			return errors.Wrap(err, "Failed to copy running state to candidate state")
		}

		s.RunningState = newState
		return s.store.Save(ctx, s.RunningState.Hosts)
	})
}

func (s *BeatStateService) copyRunningStateWithoutDevice(deviceId string) (*BeatState, error) {
	// Copy the running config to the a new config and filter any existing version of this device
	var newState = new(BeatState)
	var newHosts = make([]Host, len(s.RunningState.Hosts))
	var n = 0
	for i, h := range s.RunningState.Hosts {
		if h.Id != deviceId {
			newHosts[i-n] = h
		} else {
			n = n + 1
		}
	}
	newState.Hosts = newHosts[:len(newHosts)-n]
	return newState, nil
}

func mapDeviceStates(d []api.DeviceDTO) (*BeatState, error) {
	var c BeatState
	for _, v := range d {
		host, err := mapDeviceState(v)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to map device state")
		}
		c.Hosts = append(c.Hosts, host)
	}
	return &c, nil
}

func mapDeviceState(v api.DeviceDTO) (h Host, err error) {
	h.Id = v.DeviceId
	h.ServiceId = v.ServiceId
	h.Host = v.Ip
	h.Tags = v.Tags
	return h, nil
}

func unmarshalState(data []byte) (*BeatState, error) {
	var h []Host
	if data == nil {
		data = []byte("[]")
	}
	err := json.Unmarshal(data, &h)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize hosts")
	}

	return &BeatState{
		Hosts: h,
	}, nil
}

func newStateService(ctx context.Context) (*BeatStateService, error) {
	store, err := state.NewStateStore(ctx)
	if err != nil {
		return nil, err
	}

	result := &BeatStateService{
		store:  store,
		worker: types.NewWorker(ctx),
	}

	webconfig.RegisterService(result)

	return result, nil
}
