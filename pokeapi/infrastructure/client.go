package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	Results []ListPokemonResponseResult `json:"results"`
}

type ListPokemonResponseResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type DetailPokemonResponse struct {
	ID       uint64                       `json:"id"`
	Name     string                       `json:"name"`
	Order    int64                        `json:"order"`
	Height   uint64                       `json:"height"`
	Weight   uint64                       `json:"weight"`
	Stats    []DetailPokemonResponseStats `json:"stats"`
	Types    []DetailPokemonResponseTypes `json:"types"`
	Sprites  DetailPokemonResponseSprites `json:"sprites"`
	Versions map[string]interface{}       `json:"versions"`
}

type DetailPokemonResponseStats struct {
	BaseStat uint64                    `json:"base_stat"`
	Stat     DetailPokemonResponseStat `json:"stat"`
}

type DetailPokemonResponseStat struct {
	Name string `json:"name"`
}

type DetailPokemonResponseTypes struct {
	Slot uint64                    `json:"slot"`
	Type DetailPokemonResponseType `json:"type"`
}

type DetailPokemonResponseType struct {
	Name string `json:"name"`
}

type DetailPokemonResponseSprites struct {
	FrontDefault string                            `json:"front_default"`
	Other        DetailPokemonResponseSpritesOther `json:"other"`
}

type DetailPokemonResponseSpritesOther struct {
	OfficialArtwork DetailPokemonResponseSpritesOtherOfficialArtwork `json:"official-artwork"`
}

type DetailPokemonResponseSpritesOtherOfficialArtwork struct {
	FrontDefault string `json:"front_default"`
}

func NewClient(conn Connection) Client {
	return Client{
		conn:    conn,
		baseUrl: "https://pokeapi.co/api/v2",
	}
}

func (c *Client) ListPokemon(offset uint, limit uint) (ListPokemonResponse, error) {
	var data ListPokemonResponse

	query := map[string]string{
		"offset": fmt.Sprintf("%v", offset),
		"limit":  fmt.Sprintf("%v", limit),
	}

	u, err := c.parseUrl("/pokemon", query)
	if err != nil {
		return data, err
	}

	resp, err := c.conn.Get(u)
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

	return data, err
}

func (c *Client) GetPokemon(id uint) (DetailPokemonResponse, error) {
	var data DetailPokemonResponse

	u, err := c.parseUrl(fmt.Sprintf("/pokemon/%v", id), map[string]string{})
	if err != nil {
		return data, err
	}

	resp, err := c.conn.Get(u)
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

func (c *Client) parseUrl(path string, query map[string]string) (string, error) {
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
