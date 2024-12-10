package remote

import (
	"fmt"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client"
	"github.com/apex/log"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
)

type Connector struct {
	Config config.Api

	token string
}

func (c *Connector) ObtainToken() error {
	ok := false
	c.token, ok = tryReadToken()
	if ok {
		log.Infof("API Token found: %s", c.token)
		return nil
	}

	if err := c.signUp(); err != nil {
		return fmt.Errorf("registration on API server failed: %w", err)
	}
	log.Infof("API Token obtained: %s", c.token)

	if !tryWriteToken(c.token) {
		log.Warn("Cannot write API Token to the filesystem. Some features will be unstable")
		return nil
	}

	return nil
}

func (c *Connector) NewDiscoveryClient(apiPath string) (auth runtime.ClientAuthInfoWriter, cli *client.Client) {
	tr := httptransport.New(c.Config.Host, apiPath, []string{c.Config.Scheme})
	auth = httptransport.APIKeyAuth("X-Token", "header", c.token)
	cli = client.New(tr, strfmt.Default)
	return
}
