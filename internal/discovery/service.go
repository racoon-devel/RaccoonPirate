package discovery

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client"
	"github.com/go-openapi/runtime"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
)

const searchResultsLimit = 8

type Service struct {
	auth runtime.ClientAuthInfoWriter
	cli  *client.Client
}

type ClientFactory interface {
	NewDiscoveryClient(apiPath string) (auth runtime.ClientAuthInfoWriter, cli *client.Client)
}

func NewService(clientFactory ClientFactory, conf config.Discovery) *Service {
	s := Service{}
	s.auth, s.cli = clientFactory.NewDiscoveryClient(conf.ApiPath)
	return &s
}
