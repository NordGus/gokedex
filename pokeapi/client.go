package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Connection interface {
	Get(string) (*http.Response, error)
	CloseIdleConnections()
}

type Client struct {
	conn    Connection
	baseUrl string
}

type ListPokemonResponse struct {
	Count   uint64                      `json:"count"`
	Next    string                      `json:"next"`
	Prev    string                      `json:"prev"`
	Results []ListPokemonResponseResult `json:"results"`
}

type ListPokemonResponseResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func NewClient(conn Connection) Client {
	return Client{
		conn:    conn,
		baseUrl: "https://pokeapi.co/api/v2",
	}
}

func (c *Client) ListPokemon(offset uint, limit uint) (data ListPokemonResponse, err error) {
	url := fmt.Sprintf("%s/pokemon", c.baseUrl)

	resp, err := c.conn.Get(url)
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
