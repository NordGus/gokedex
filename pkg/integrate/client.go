package integrate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (c *client) createPokemonPage(pokemon PokemonPage) (NotionPageCreatedResponse, error) {
	var data NotionPageCreatedResponse

	reqUrl, err := c.parseUrl("/pages", map[string]string{})
	if err != nil {
		return data, err
	}

	postBody, err := json.Marshal(pokemon)
	if err != nil {
		return data, err
	}

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(postBody))
	if err != nil {
		return data, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.secret))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	resp, err := c.conn.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

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
