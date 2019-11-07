package loginterceptor

import (
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"encoding/json"
	"net/http"
)

var (
	logger = log.NewLogger("msx.httpclient.loginterceptor")
)

func NewInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	return func(req *http.Request) (response *http.Response, err error) {
		ctx := req.Context()
		response, err = fn(req)
		if response == nil {
			logger.WithContext(ctx).WithError(err).Error("000 : ", req.URL.String())
		} else if response.StatusCode > 399 {
			// Fully log the response
			logger.WithContext(ctx).Errorf("%s : %s", response.Status, req.URL.String())
			var responseBytes []byte
			responseBytes, _ = json.Marshal(response)
			logger.WithContext(ctx).Error(string(responseBytes))
		} else {
			logger.WithContext(ctx).Infof("%s : %s", response.Status, req.URL.String())
		}

		return response, err
	}
}
