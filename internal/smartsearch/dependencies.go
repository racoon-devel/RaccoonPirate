package smartsearch

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client"
	"github.com/go-openapi/runtime"
)

type ClientFactory interface {
	NewDiscoveryClient(apiPath string) (auth runtime.ClientAuthInfoWriter, cli *client.Client)
}
