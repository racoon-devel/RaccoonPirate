package remote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type signUpResponse struct {
	Token string
}

func (c *Connector) signUp() error {
	u := url.URL{}

	query := url.Values{}
	query.Add("domain", c.Config.Domain)

	u.Scheme = c.Config.Scheme
	u.Host = fmt.Sprintf("%s:%d", c.Config.Host, c.Config.Port)
	u.RawQuery = query.Encode()
	u.Path = "/signup"

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d [ %s ]", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	tokenResponse := signUpResponse{}
	if err = json.Unmarshal(body, &tokenResponse); err != nil {
		return err
	}

	c.token = tokenResponse.Token
	return nil
}
