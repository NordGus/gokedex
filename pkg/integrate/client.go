package integrate

import (
	"fmt"
	"net/url"
)

type IntegrationSecret string

type DatabaseID string

type client struct {
	conn    Connection
	baseUrl string
	secret  IntegrationSecret
}

func newClient(conn Connection, secret IntegrationSecret) client {
	return client{
		conn:    conn,
		baseUrl: "https://api.notion.com/v1",
		secret:  secret,
	}
}

// func (c *client) createPokemonPage(pokemon Pokemon) (PokemonResponse, error) {
// 	req := http.NewRequest("GET", c.parseUrl())
// }

func (c *client) parseUrl(path string, query map[string]string) (string, error) {
	u, err := url.Parse(fmt.Sprintf("%v%v", c.baseUrl, path))
	if err != nil {
		return "", err
	}

	q := u.Query()

	for k, v := range query {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}
