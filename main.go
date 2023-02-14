package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/NordGus/gokedex/pokeapi"
)

func main() {
	client := http.Client{
		Timeout: time.Millisecond * 3000,
	}

	apiClient := pokeapi.NewClient(&client)
	limit := 20
	output := make(chan pokeapi.DetailPokemonResponse)
	resultCount := 0

	go func(client pokeapi.Client, output chan pokeapi.DetailPokemonResponse, limit int) {
		var wg sync.WaitGroup

		for offset := 0; ; offset += limit {
			data, err := client.ListPokemon(uint(offset), uint(limit))
			if err != nil {
				panic(err.Error())
			}

			resultCount = int(data.Count)

			for _, result := range data.Results {
				wg.Add(1)
				go getPokemon(&wg, client, result, output)
			}

			if data.Next == "" {
				wg.Wait()
				close(output)
				break
			}
		}
	}(apiClient, output, limit)

	count := 0

	for range output {
		count++
	}

	fmt.Println(count)
	fmt.Println(resultCount)
	fmt.Println(count == resultCount)
}

func getPokemon(wg *sync.WaitGroup, client pokeapi.Client, data pokeapi.ListPokemonResponseResult, out chan<- pokeapi.DetailPokemonResponse) {
	defer wg.Done()

	u, err := url.Parse(data.Url)
	if err != nil {
		panic(err.Error())
	}

	path := []string{}

	for _, s := range strings.Split(u.Path, "/") {
		if s != "" {
			path = append(path, s)
		}
	}

	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		panic(err.Error())
	}

	result, err := client.GetPokemon(uint(id))
	if err != nil {
		panic(err.Error())
	}

	out <- result
}
