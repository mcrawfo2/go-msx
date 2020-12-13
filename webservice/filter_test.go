package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"github.com/emicklei/go-restful"
)

const HeaderNameDummyFilter = "X-Dummy-Filter"
const HeaderValueDummyFilter = "true"

var DummyFilterResponseCheck = webservicetest.ResponseHasHeader(HeaderNameDummyFilter, HeaderValueDummyFilter)

func DummyFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	defer func() {
		resp.Header().Set(HeaderNameDummyFilter, HeaderValueDummyFilter)
		chain.ProcessFilter(req, resp)
	}()
}

