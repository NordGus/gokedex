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

func (s *Service) IntegrateToNotion(in <-chan extract.Pokemon, limits chan bool) <-chan struct{} {
	pages := s.preparePokemonPages(in, limits)
	responses := s.createPokedexPages(pages, limits)
	done := s.logPageCreated(responses, limits)

	return done
}

func (s *Service) preparePokemonPages(in <-chan extract.Pokemon, limits chan bool) <-chan PokemonPage {
	var wg sync.WaitGroup
	out := make(chan PokemonPage)

	go func(wg *sync.WaitGroup, in <-chan extract.Pokemon, out chan<- PokemonPage, limits chan bool) {
		defer close(out)
		defer func(limits <-chan bool) {
			<-limits
		}(limits)

		limits <- true

		for pokemon := range in {
			wg.Add(1)
			go s.mapPokemonPage(wg, pokemon, out, limits)
		}

		wg.Wait()
	}(&wg, in, out, limits)

	return out
}

func (s *Service) mapPokemonPage(wg *sync.WaitGroup, pokemon extract.Pokemon, out chan<- PokemonPage, limits chan bool) {
	defer wg.Done()
	defer func(limits <-chan bool) {
		<-limits
	}(limits)

	limits <- true

	out <- externalPokemonToInternalPokemon(pokemon, s.databaseID)
}

func (s *Service) createPokedexPages(in <-chan PokemonPage, limits chan bool) <-chan NotionPageCreatedResponse {
	var wg sync.WaitGroup
	out := make(chan NotionPageCreatedResponse)

	go func(wg *sync.WaitGroup, in <-chan PokemonPage, out chan<- NotionPageCreatedResponse, limits chan bool) {
		defer close(out)
		defer func(limits <-chan bool) {
			<-limits
		}(limits)

		limits <- true

		for page := range in {
			wg.Add(1)
			go s.createPage(wg, page, out, limits)
		}

		wg.Wait()
	}(&wg, in, out, limits)

	return out
}

func (s *Service) createPage(wg *sync.WaitGroup, page PokemonPage, out chan<- NotionPageCreatedResponse, limits chan bool) {
	defer wg.Done()
	defer func(limits <-chan bool) {
		<-limits
	}(limits)

	limits <- true

	resp, err := s.client.createPokemonPage(page)
	if err != nil {
		panic(err.Error())
	}

	out <- resp
}

func (s *Service) logPageCreated(in <-chan NotionPageCreatedResponse, limits chan bool) <-chan struct{} {
	out := make(chan struct{})

	go func(in <-chan NotionPageCreatedResponse, out chan<- struct{}, limits chan bool) {
		defer close(out)
		defer func(limits <-chan bool) {
			<-limits
		}(limits)

		limits <- true

		processed := 0

		for page := range in {
			processed++
			fmt.Println("Pekédex Page Created:", page.Url)
		}

		fmt.Println("Processed Pokémon:", processed)
	}(in, out, limits)

	return out
}
