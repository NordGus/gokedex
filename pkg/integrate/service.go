package integrate

import (
	"fmt"
	"sync"

	"github.com/NordGus/gokedex/pkg/extract"
)

type Service struct {
	client     client
	databaseID DatabaseID
	sem        chan bool
}

func NewService(sem chan bool, conn Connection, connSem chan bool, secret IntegrationSecret, databaseId DatabaseID) Service {
	return Service{
		client:     newClient(conn, connSem, secret),
		databaseID: databaseId,
		sem:        sem,
	}
}

func (s *Service) IntegrateToNotion(in <-chan extract.Pokemon) <-chan struct{} {
	pages := s.preparePokemonPages(in)
	responses := s.createPokedexPages(pages)
	done := s.logPageCreated(responses)

	return done
}

func (s *Service) freeResources() {
	<-s.sem
}

func (s *Service) preparePokemonPages(in <-chan extract.Pokemon) <-chan PokemonPage {
	var wg sync.WaitGroup
	out := make(chan PokemonPage)

	go func(wg *sync.WaitGroup, in <-chan extract.Pokemon, out chan<- PokemonPage) {
		s.sem <- true
		defer s.freeResources()
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
	s.sem <- true
	defer s.freeResources()
	defer wg.Done()

	out <- externalPokemonToInternalPokemon(pokemon, s.databaseID)
}

func (s *Service) createPokedexPages(in <-chan PokemonPage) <-chan NotionPageCreatedResponse {
	var wg sync.WaitGroup
	out := make(chan NotionPageCreatedResponse)

	go func(wg *sync.WaitGroup, in <-chan PokemonPage, out chan<- NotionPageCreatedResponse) {
		s.sem <- true
		defer s.freeResources()
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
	s.sem <- true
	defer s.freeResources()
	defer wg.Done()

	resp, err := s.client.createPokemonPage(page)
	if err != nil {
		panic(err.Error())
	}

	out <- resp
}

func (s *Service) logPageCreated(in <-chan NotionPageCreatedResponse) <-chan struct{} {
	var wg sync.WaitGroup

	out := make(chan struct{})

	go func(wg *sync.WaitGroup, in <-chan NotionPageCreatedResponse, out chan<- struct{}) {
		s.sem <- true
		defer s.freeResources()
		defer close(out)

		processed := 0

		for page := range in {
			wg.Add(1)
			processed++
			fmt.Println("Pekédex Page Created:", page.Url)
			wg.Done()
		}

		wg.Wait()

		fmt.Println("Processed Pokémon:", processed)
	}(&wg, in, out)

	return out
}
