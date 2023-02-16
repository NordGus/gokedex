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
	pokemon := s.preparePokemonPages(in)
	pages := s.createPokedexPages(pokemon)
	done := s.logPageCreated(pages)

	return done
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

func (s *Service) createPokedexPages(in <-chan PokemonPage) <-chan NotionPageCreatedResponse {
	var wg sync.WaitGroup
	out := make(chan NotionPageCreatedResponse)

	go func(wg *sync.WaitGroup, in <-chan PokemonPage, out chan<- NotionPageCreatedResponse) {
		defer close(out)

		for page := range in {
			wg.Add(1)
			go s.createPage(wg, page, out)
		}

		wg.Wait()
	}(&wg, in, out)

	return out
}

func (s *Service) createPage(wg *sync.WaitGroup, page PokemonPage, out chan<- NotionPageCreatedResponse) {
	defer wg.Done()

	resp, err := s.client.createPokemonPage(page)
	if err != nil {
		panic(err.Error())
	}

	out <- resp
}

func (s *Service) logPageCreated(in <-chan NotionPageCreatedResponse) <-chan struct{} {
	out := make(chan struct{})

	go func(in <-chan NotionPageCreatedResponse, out chan<- struct{}) {
		defer close(out)
		pages := []NotionPageCreatedResponse{}

		for page := range in {
			pages = append(pages, page)
			fmt.Println("Pekédex Page Created:", page.Url)
		}

		fmt.Println("Processed Pokémon:", len(pages))
	}(in, out)

	return out
}
