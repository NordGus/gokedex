package extract

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

type client struct {
	conn    Connection
	baseUrl string
	sem     chan bool
}

func newClient(conn Connection, sem chan bool) client {
	return client{
		conn:    conn,
		baseUrl: "https://pokeapi.co/api/v2",
		sem:     sem,
	}
}

func (c *client) freeConnection() {
	<-c.sem
}

func (c *client) listPokemonSpecies(offset uint, limit uint) (pokemonSpeciesPageResponse, error) {
	c.sem <- true
	defer c.freeConnection()

	var data pokemonSpeciesPageResponse

	query := map[string]string{
		"offset": fmt.Sprintf("%v", offset),
		"limit":  fmt.Sprintf("%v", limit),
	}

	u, err := c.parseUrl("/pokemon-species", query)
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

func (c *client) getPokemonSpecies(id uint) (pokemonSpeciesResponse, error) {
	c.sem <- true
	defer c.freeConnection()

	var data pokemonSpeciesResponse

	u, err := c.parseUrl(fmt.Sprintf("/pokemon-species/%v", id), map[string]string{})
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

func (c *client) getPokemon(id uint) (pokemonResponse, error) {
	c.sem <- true
	defer c.freeConnection()

	var data pokemonResponse

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
