package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NordGus/gokedex/pokeapi"
)

func main() {
	client := http.Client{
		Timeout: time.Millisecond * 3000,
	}

	apiClient := pokeapi.NewClient(&client)

	data, err := apiClient.ListPokemon(uint(0), uint(20))
	if err != nil {
		panic(err.Error())
	}

	for _, result := range data.Results {
		fmt.Printf("\tName: %v,\tUrl: %v\n", result.Name, result.Url)
	}
}
