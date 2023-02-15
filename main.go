package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/NordGus/gokedex/pokeapi"
)

func main() {
	client := http.Client{
		Timeout: time.Millisecond * 3000,
	}

	service := pokeapi.NewService(&client)
	start := time.Now()
	pokemon := service.ExtractPokemon()

	fmt.Println("Duration:", time.Since(start))

	different := make(map[int64][]pokeapi.Pokemon)

	for _, pk := range pokemon {
		different[pk.Number] = append(different[pk.Number], pk)
	}

	random := rand.New(rand.NewSource(time.Now().Unix()))

	fmt.Println("Different Pokémon:", len(different))
	fmt.Println("Repited Pokémon:", (len(pokemon) - len(different)))
	fmt.Println("Processed Pokémon:", len(pokemon))
	fmt.Printf("Processed Pokémon Example:\n%+v\n", pokemon[random.Intn(len(pokemon))])
}
