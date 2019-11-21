package dnac

import "cto-github.cisco.com/NFV-BU/go-msx/integration"

type Api interface {
	Connect(request DnacConnectRequest) (*integration.MsxResponse, error)
	RetrieveExtendedData(request DnacExtendedRequest) (*integration.MsxResponse, error)
}
