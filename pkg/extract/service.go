package extract

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

const (
	listPokemonSpeciesChannelBufferSize   = 5
	getPokemonSpeciesChannelBuffersize    = 5
	getPokemonDetailsChannelBufferSize    = 5
	buildPokemonChannelBufferSize         = 5
	logPokemonExtractionChannelBufferSize = 5
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

func (s *Service) ExtractPokemon(limits chan bool) <-chan Pokemon {
	pages := s.listPokemonSpecies(limits)
	species := s.getPokemonSpecies(pages, limits)
	details := s.getPokemonDetails(species, limits)
	pokemon := s.buildPokemon(details, limits)
	results := s.logPokemonExtraction(pokemon, limits)

	return results
}

func (s *Service) listPokemonSpecies(limits chan bool) <-chan pokemonSpeciesPageResponse {
	offset := uint(0)
	limit := uint(20)
	out := make(chan pokemonSpeciesPageResponse, listPokemonSpeciesChannelBufferSize)

	go func(offset uint, limit uint, out chan<- pokemonSpeciesPageResponse, limits chan bool) {
		defer close(out)
		defer func(limits <-chan bool) {
			<-limits
		}(limits)

		limits <- true

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
	}(offset, limit, out, limits)

	return out
}

func (s *Service) getPokemonSpecies(in <-chan pokemonSpeciesPageResponse, limits chan bool) <-chan pokemonSpeciesResponse {
	var wg sync.WaitGroup

	out := make(chan pokemonSpeciesResponse, getPokemonSpeciesChannelBuffersize)

	go func(wg *sync.WaitGroup, in <-chan pokemonSpeciesPageResponse, out chan<- pokemonSpeciesResponse, limits chan bool) {
		defer close(out)
		defer func(limits <-chan bool) {
			<-limits
		}(limits)

		limits <- true

		for resp := range in {
			wg.Add(len(resp.Results))

			for _, species := range resp.Results {
				go s.retrievePokemonSpecies(wg, parseIdFromResponseUrl(species.Url), out, limits)
			}
		}

		wg.Wait()
	}(&wg, in, out, limits)

	return out
}

func (s *Service) retrievePokemonSpecies(wg *sync.WaitGroup, id uint, out chan<- pokemonSpeciesResponse, limits chan bool) {
	defer wg.Done()
	defer func(limits <-chan bool) {
		<-limits
	}(limits)

	limits <- true

	data, err := s.client.getPokemonSpecies(uint(id))
	if err != nil {
		panic(err.Error())
	}

	out <- data
}

func (s *Service) getPokemonDetails(in <-chan pokemonSpeciesResponse, limits chan bool) <-chan fullPokemonData {
	var wg sync.WaitGroup

	out := make(chan fullPokemonData, getPokemonDetailsChannelBufferSize)

	go func(wg *sync.WaitGroup, in <-chan pokemonSpeciesResponse, out chan<- fullPokemonData, limits chan bool) {
		defer close(out)
		defer func(limits <-chan bool) {
			<-limits
		}(limits)

		limits <- true

		for resp := range in {
			wg.Add(1)

			go s.retrievePokemonDetail(wg, resp, out, limits)
		}

		wg.Wait()
	}(&wg, in, out, limits)

	return out
}

func (s *Service) retrievePokemonDetail(wg *sync.WaitGroup, species pokemonSpeciesResponse, out chan<- fullPokemonData, limits chan bool) {
	defer wg.Done()
	defer func(limits <-chan bool) {
		<-limits
	}(limits)

	limits <- true

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

func (s *Service) buildPokemon(in <-chan fullPokemonData, limits chan bool) <-chan Pokemon {
	var wg sync.WaitGroup

	out := make(chan Pokemon, buildPokemonChannelBufferSize)

	go func(wg *sync.WaitGroup, in <-chan fullPokemonData, out chan<- Pokemon, limits chan bool) {
		defer close(out)
		defer func(limits <-chan bool) {
			<-limits
		}(limits)

		limits <- true

		for resp := range in {
			wg.Add(1)
			go s.mapDataToPokemon(wg, resp, out, limits)
		}

		wg.Wait()
	}(&wg, in, out, limits)

	return out
}

func (s *Service) mapDataToPokemon(wg *sync.WaitGroup, data fullPokemonData, out chan<- Pokemon, limits chan bool) {
	defer wg.Done()
	defer func(limits <-chan bool) {
		<-limits
	}(limits)

	limits <- true

	pokemon := mapPokemon(data.Species, data.Details)

	out <- pokemon
}

func (s *Service) logPokemonExtraction(in <-chan Pokemon, limits chan bool) <-chan Pokemon {
	var wg sync.WaitGroup

	out := make(chan Pokemon, logPokemonExtractionChannelBufferSize)

	go func(wg *sync.WaitGroup, in <-chan Pokemon, out chan<- Pokemon, limits chan bool) {
		defer close(out)
		defer func(limits <-chan bool) {
			<-limits
		}(limits)

		limits <- true

		for pokemon := range in {
			wg.Add(1)

			fmt.Println(pokemon.Name, "extracted from PokéAPI!")
			out <- pokemon

			wg.Done()
		}

		wg.Wait()
	}(&wg, in, out, limits)

	return out
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
