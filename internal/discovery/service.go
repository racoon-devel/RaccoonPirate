package discovery

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/racoon-devel/media-station/internal/config"
	"sync"
)

const searchResultsLimit = 8

type Service struct {
	auth  runtime.ClientAuthInfoWriter
	cli   *client.Client
	cache sync.Map
}

func NewService(conf config.Discovery) *Service {
	tr := httptransport.New(conf.Host, conf.Path, []string{conf.Scheme})
	auth := httptransport.APIKeyAuth("X-Token", "header", conf.Identity)
	cli := client.New(tr, strfmt.Default)

	return &Service{
		auth: auth,
		cli:  cli,
	}
}
