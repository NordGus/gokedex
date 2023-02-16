package extract

import (
	"strconv"
	"strings"
	"sync"
)

type Service struct {
	client client
}

type fullPokemonData struct {
	Species pokemonSpeciesResponse
	Details pokemonResponse
}

func NewService(conn Connection) Service {
	return Service{
		client: newClient(conn),
	}
}

func (s *Service) ExtractPokemon() <-chan Pokemon {
	pages := s.listPokemonSpecies()
	species := s.getPokemonSpecies(pages)
	details := s.getPokemonDetails(species)
	pokemon := s.buildPokemon(details)

	return pokemon
}

func (s *Service) listPokemonSpecies() <-chan pokemonSpeciesPageResponse {
	offset := uint(0)
	limit := uint(20)
	out := make(chan pokemonSpeciesPageResponse)

	go func(offset uint, limit uint, out chan<- pokemonSpeciesPageResponse) {
		defer close(out)

		for ; ; offset += limit {
			data, err := s.client.listPokemonSpecies(offset, limit)
			if err != nil {
				panic(err.Error())
			}

			out <- data

			if data.Next == "" {
				break
			}
		}
	}(offset, limit, out)

	return out
}

func (s *Service) getPokemonSpecies(in <-chan pokemonSpeciesPageResponse) <-chan pokemonSpeciesResponse {
	var wg sync.WaitGroup

	out := make(chan pokemonSpeciesResponse)

	go func(wg *sync.WaitGroup, in <-chan pokemonSpeciesPageResponse, out chan<- pokemonSpeciesResponse) {
		defer close(out)

		for resp := range in {
			wg.Add(len(resp.Results))

			for _, species := range resp.Results {
				go s.retrievePokemonSpecies(wg, parseIdFromResponseUrl(species.Url), out)
			}
		}

		wg.Wait()
	}(&wg, in, out)

	return out
}

func (s *Service) retrievePokemonSpecies(wg *sync.WaitGroup, id uint, out chan<- pokemonSpeciesResponse) {
	defer wg.Done()

	data, err := s.client.getPokemonSpecies(uint(id))
	if err != nil {
		panic(err.Error())
	}

	out <- data
}

func (s *Service) getPokemonDetails(in <-chan pokemonSpeciesResponse) <-chan fullPokemonData {
	var wg sync.WaitGroup

	out := make(chan fullPokemonData)

	go func(wg *sync.WaitGroup, in <-chan pokemonSpeciesResponse, out chan<- fullPokemonData) {
		defer close(out)

		for resp := range in {
			wg.Add(1)

			go s.retrievePokemonDetail(wg, resp, out)
		}

		wg.Wait()
	}(&wg, in, out)

	return out
}

func (s *Service) retrievePokemonDetail(wg *sync.WaitGroup, species pokemonSpeciesResponse, out chan<- fullPokemonData) {
	defer wg.Done()

	for _, variety := range species.Varieties {
		if variety.IsDefault {
			data, err := s.client.getPokemon(parseIdFromResponseUrl(variety.Pokemon.Url))
			if err != nil {
				panic(err.Error())
			}

			out <- fullPokemonData{
				Species: species,
				Details: data,
			}

			return
		}
	}
}

func (s *Service) buildPokemon(in <-chan fullPokemonData) <-chan Pokemon {
	var wg sync.WaitGroup

	out := make(chan Pokemon)

	go func(wg *sync.WaitGroup, in <-chan fullPokemonData, out chan<- Pokemon) {
		defer close(out)

		for resp := range in {
			wg.Add(1)
			go s.mapDataToPokemon(wg, resp, out)
		}

		wg.Wait()
	}(&wg, in, out)

	return out
}

func (s *Service) mapDataToPokemon(wg *sync.WaitGroup, data fullPokemonData, out chan<- Pokemon) {
	defer wg.Done()

	out <- mapPokemon(data.Species, data.Details)
}

func parseIdFromResponseUrl(respUrl string) uint {
	comp := strings.Split(respUrl, "/")

	for i := len(comp) - 1; i >= 0; i-- {
		if comp[i] != "" {
			id, err := strconv.Atoi(comp[i])
			if err != nil {
				return 0
			}

			return uint(id)
		}
	}

	return 0
}
