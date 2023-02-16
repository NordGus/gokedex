package integrate

import (
	"fmt"
	"sync"

	"github.com/NordGus/gokedex/pkg/extract"
)

type Service struct {
	client     client
	databaseID DatabaseID
}

func NewService(conn Connection, secret IntegrationSecret, databaseId DatabaseID) Service {
	return Service{
		client:     newClient(conn, secret),
		databaseID: databaseId,
	}
}

func (s *Service) IntegrateToNotion(in <-chan extract.Pokemon) <-chan struct{} {
	out := make(chan struct{})
	pokemon := s.preparePokemonPages(in)

	go func(in <-chan PokemonPage, out chan<- struct{}) {
		defer close(out)
		pokemon := []PokemonPage{}

		for poke := range in {
			pokemon = append(pokemon, poke)
			fmt.Printf("Pokemon Page:\n%+v\n", poke)
		}

		fmt.Println("Processed Pokémon:", len(pokemon))
	}(pokemon, out)

	return out
}

func (s *Service) preparePokemonPages(in <-chan extract.Pokemon) <-chan PokemonPage {
	var wg sync.WaitGroup
	out := make(chan PokemonPage)

	go func(wg *sync.WaitGroup, in <-chan extract.Pokemon, out chan<- PokemonPage) {
		defer close(out)

		for pokemon := range in {
			wg.Add(1)
			go s.mapPokemonPage(wg, pokemon, out)
		}

		wg.Wait()
	}(&wg, in, out)

	return out
}

func (s *Service) mapPokemonPage(wg *sync.WaitGroup, pokemon extract.Pokemon, out chan<- PokemonPage) {
	defer wg.Done()

	out <- externalPokemonToInternalPokemon(pokemon, s.databaseID)
}
