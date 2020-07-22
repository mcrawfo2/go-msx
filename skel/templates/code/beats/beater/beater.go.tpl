package beater

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx-beats/locker"
	"cto-github.cisco.com/NFV-BU/go-msx-beats/meta"
	"cto-github.cisco.com/NFV-BU/go-msx-beats/publisher"
	"cto-github.cisco.com/NFV-BU/go-msx-beats/worker"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/elastic/beats/libbeat/common"
	"github.com/pkg/errors"
	"sort"
	"time"
)

const (
	FieldUs       = "us"
	FieldDuration = "duration"
)

var logger = log.NewLogger("${app.name}.internal.beater")

type Beater struct {
	config       *BeatConfig
	stateService *BeatStateService
	locks        *locker.Locker
}

// Heartbeat main loop
func (b *Beater) Run(ctx context.Context) {
	logger.WithContext(ctx).Info("Running heartbeat loop")

	ticker := time.NewTicker(b.config.Period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		batches := b.batches(ctx, b.stateService.RunningState.Hosts)
		b.scheduleBatches(ctx, batches)
	}
}

func (b *Beater) batches(ctx context.Context, allHosts []Host) [][]Host {
	batches := make([][]Host, 0)
	batchIdx := -1
	batchSize := b.config.Batch.Size
	if batchSize == 0 {
		batchSize = len(allHosts)
	}

	hosts := make([]Host, 0)
	for _, h := range allHosts {
		hosts = append(hosts, h)
	}
	sort.SliceStable(hosts, func(i, j int) bool { return hosts[i].Id < hosts[j].Id })

	for _, h := range hosts {
		logger.WithContext(ctx).Debugf("Scheduling host: %s", h.Host)

		if len(batches) > batchIdx || len(batches[batchIdx]) == batchSize {
			batchIdx = batchIdx + 1
			batches = append(batches, make([]Host, 0))
		}
		batches[batchIdx] = append(batches[batchIdx], h)
	}

	return batches
}

func (b *Beater) scheduleBatches(ctx context.Context, batches [][]Host) {
	now := time.Now()
	stride := b.config.Batch.Delay
	scheduler := worker.NewScheduler()
	for n := range batches {
		logger.WithContext(ctx).Debugf("Scheduling batch: %d", n)

		batch := batches[n]
		for _, h := range batch {
			scheduler.Schedule(ctx, now, b.hostAction(h))
		}

		now = now.Add(stride)
	}
}

func (b *Beater) hostAction(h Host) types.ActionFunc {
	return func(ctx context.Context) error {
		defer worker.Recovery(logger)

		ctx, span := trace.NewSpan(
			trace.UntracedContextFromContext(ctx),
			"probe")
		span.SetTag(trace.FieldDeviceId, h.Id)
		span.SetTag(trace.FieldServiceId, h.ServiceId)
		span.SetTag(trace.FieldDeviceAddress, h.Host)
		defer span.Finish()

		// Ensure previous tasks have completed for this host
		ok := b.locks.TestAndClaim(h.Host)
		if !ok {
			b.locks.Timeout(h.Host, 4)
			return nil
		}
		defer b.locks.Release(h.Host)

		if err := b.host(ctx, h); err != nil {
			span.SetTag(trace.FieldError, err.Error())
		}

		return nil
	}
}

func (b *Beater) host(ctx context.Context, h Host) (err error) {
	logFields := h.NewLog()
	logger.WithContext(ctx).
		WithFields(logFields).
		Infof("Retrieving metrics for device %q", h.Id)

	startTime := time.Now()

	var result = make(common.MapStr)

	// Probe the host for a document
	result, err = b.dummy(h)

	// Apply the results to the event document
	fields := h.NewEvent()
	fields.Update(result)

	// Apply failure to the event document if an error occurred
	failed := err != nil
	if failed {
		logger.
			WithContext(ctx).
			WithError(err).
			WithFields(logFields).
			Errorf("Failed to retrieve metrics for device %q", h.Id)
		fields[meta.FieldFailed] = true
		fields.Put(meta.FieldErrorMessage, err.Error())
	}

	// Apply action duration to the event document
	actionDuration := time.Now().Sub(startTime)
	fields[FieldDuration] = valueMap(FieldUs, int64(actionDuration/time.Microsecond))

	if publisherErr := publisher.Publish(ctx, startTime, fields); publisherErr != nil {
		logger.
			WithContext(ctx).
			WithError(publisherErr).
			WithFields(logFields).
			Error("Failed to publish event")
		if err == nil {
			err = publisherErr
		}
	}

	return
}

func valueMap(unit string, value interface{}) common.MapStr {
	return common.MapStr{
		unit: value,
	}
}

func (b *Beater) dummy(host Host) (results common.MapStr, err error) {
	return common.MapStr{
		meta.FieldUp: true,
	}, nil
}

func newBeater(ctx context.Context) (*Beater, error) {
	beatConfig, err := newConfig(config.FromContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create config")
	}

	stateService, err := newStateService(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create state service")
	}

	return &Beater{
		config:       beatConfig,
		stateService: stateService,
		locks:        locker.NewLocker(),
	}, nil
}
