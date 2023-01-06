// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package metricsprovider

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cache/lru"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/repository"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	promapi "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	endpointName = "metrics"
)

type MetricTags map[string]string

func (t MetricTags) SetTag(name, value string) {
	t[name] = value
}

func (t MetricTags) MatchesMetric(metric *promapi.Metric) bool {
	if len(t) == 0 {
		return true
	}

	for _, label := range metric.GetLabel() {
		tagValue := t[label.GetName()]
		if tagValue != label.GetValue() {
			return false
		}
	}

	return true
}

func (t MetricTags) MatchingMetrics(family *promapi.MetricFamily) []*promapi.Metric {
	var results []*promapi.Metric
	for _, metric := range family.GetMetric() {
		if t.MatchesMetric(metric) {
			results = append(results, metric)
		}
	}
	return results
}

type MetricMeasurement struct {
	Statistic string   `json:"statistic"`
	Value     *float64 `json:"value"`
}

type MetricAvailableTag struct {
	Tag    string   `json:"tag"`
	Values []string `json:"values"`
}

type MetricDefinition struct {
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	BaseUnit      string               `json:"baseUnit"`
	Measurements  []MetricMeasurement  `json:"measurements"`
	AvailableTags []MetricAvailableTag `json:"availableTags"`
}

type Report struct {
	Names []string `json:"names"`
}

type Provider struct {
	MetricCache lru.Cache
	queryMtx    sync.Mutex
}

func (h *Provider) Query() (map[string]*promapi.MetricFamily, error) {
	h.queryMtx.Lock()
	defer h.queryMtx.Unlock()

	cacheResult, ok := h.MetricCache.Get("metricFamilies")
	if ok {
		return cacheResult.(map[string]*promapi.MetricFamily), nil
	}

	var bodyBytes = bytes.NewBuffer(make([]byte, 0, 4096))

	err := push.New("localhost:50000", "query-names").
		Gatherer(prometheus.DefaultGatherer).
		Format(expfmt.FmtText).
		Client(httpclient.ClientFunc(func(req *http.Request) (*http.Response, error) {
			defer req.Body.Close()
			_, err := io.Copy(bodyBytes, req.Body)
			return &http.Response{
				StatusCode: 200,
				Body:       http.NoBody,
			}, err
		})).
		Push()
	if err != nil {
		return nil, err
	}

	mfs, err := new(expfmt.TextParser).TextToMetricFamilies(bodyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode metrics report")
	}

	h.MetricCache.Set("metricFamilies", mfs)

	return mfs, nil
}

func (h *Provider) QueryMetric(name string, tags MetricTags) (def MetricDefinition, err error) {
	mfs, err := h.Query()
	if err != nil {
		return
	}

	mf, ok := mfs[name]
	if !ok {
		err = repository.ErrNotFound
	}

	def.Name = strings.ReplaceAll(mf.GetName(), "_", ".")
	def.Description = mf.GetHelp()
	def.BaseUnit = ""

	var labels = make(map[string]types.StringSet)
	for _, metric := range mf.GetMetric() {
		for _, label := range metric.Label {
			if _, ok = labels[label.GetName()]; !ok {
				labels[label.GetName()] = make(types.StringSet)
			}
			labels[label.GetName()].Add(label.GetValue())
		}
	}

	for label, values := range labels {
		at := MetricAvailableTag{
			Tag:    label,
			Values: values.Values(),
		}
		def.AvailableTags = append(def.AvailableTags, at)
	}

	for _, metric := range tags.MatchingMetrics(mf) {
		switch {
		case metric.Gauge != nil:
			def.Measurements = append(def.Measurements, MetricMeasurement{
				Statistic: "VALUE",
				Value:     metric.Gauge.Value,
			})

		case metric.Counter != nil:
			def.Measurements = append(def.Measurements, MetricMeasurement{
				Statistic: "VALUE",
				Value:     metric.Counter.Value,
			})

		case metric.Untyped != nil:
			def.Measurements = append(def.Measurements, MetricMeasurement{
				Statistic: "VALUE",
				Value:     metric.Untyped.Value,
			})
		}
	}

	return def, nil
}

func (h *Provider) QueryMetricNames() ([]string, error) {
	mfs, err := h.Query()
	if err != nil {
		return nil, err
	}

	var results = []string{}
	for k := range mfs {
		name := strings.ReplaceAll(k, "_", ".")
		results = append(results, name)
	}

	sort.Strings(results)

	return results, nil
}

func (h *Provider) ReportMetricNames(req *restful.Request) (interface{}, error) {
	metricNames, err := h.QueryMetricNames()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query metric names")
	}

	return Report{Names: metricNames}, nil
}

func (h *Provider) ParseMetricTags(tagQueryValues []string) MetricTags {
	var metricTags = MetricTags{}
	for _, metricTag := range tagQueryValues {
		tagParts := strings.SplitN(metricTag, ":", 2)
		if len(tagParts) != 2 {
			continue
		}
		metricTags.SetTag(tagParts[0], tagParts[1])
	}
	return metricTags
}

func (h *Provider) ReportMetric(req *restful.Request) (body interface{}, err error) {
	metricName := req.PathParameter("metricName")
	metricName = strings.ReplaceAll(metricName, ".", "_")
	metricTags := h.ParseMetricTags(req.Request.URL.Query()["tag"])

	metricDef, err := h.QueryMetric(metricName, metricTags)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, webservice.NewNotFoundError(err)
		}
		return nil, webservice.NewBadRequestError(err)
	}

	return metricDef, nil
}

func (h *Provider) EndpointName() string {
	return endpointName
}

func (h *Provider) Actuate(webService *restful.WebService) error {
	webService.Consumes(restful.MIME_JSON)
	webService.Produces(restful.MIME_JSON)

	webService.Path(webService.RootPath() + "/admin/" + endpointName)

	// Unsecured routes for info
	webService.Route(webService.GET("").
		Operation("admin.metrics").
		To(adminprovider.RawAdminController(h.ReportMetricNames)).
		Do(webservice.Returns200))

	webService.Route(webService.GET("/{metricName}").
		Param(restful.QueryParameter("tag", "Specifications").AllowMultiple(true)).
		Operation("admin.metric.instance").
		To(adminprovider.RawAdminController(h.ReportMetric)).
		Do(webservice.Returns200))

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(&Provider{
			MetricCache: lru.NewCache2(
				3*time.Second,
				1000,
				10*time.Second,
				false,
				types.NewClock(ctx),
				false,
				""),
		})
		adminprovider.RegisterLink(endpointName, endpointName, false)
	}
	return nil
}
