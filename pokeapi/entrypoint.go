package pokeapi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	client "github.com/NordGus/gokedex/pokeapi/infrastructure"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service struct {
	client client.Client
}

type Pokemon struct {
	ID            uint64
	Name          string
	Number        int64
	Type          []PokemonType
	Height        uint64
	Weight        uint64
	HP            uint64
	Attack        uint64
	Defense       uint64
	SpAttack      uint64
	SpDefense     uint64
	Speed         uint64
	Category      string
	Generation    string
	Sprite        string
	Artwork       string
	BulbapediaURL string
}

type PokemonType struct {
	Name string
}

func NewService(conn client.Connection) Service {
	return Service{
		client: client.NewClient(conn),
	}
}

func (s *Service) ExtractPokemon() []Pokemon {
	listed := s.listPokemon(0, 20)
	unprocessPokemon := s.getPokemon(listed)
	pokemon := s.detailPokemonResponseToPokemon(unprocessPokemon)
	results := []Pokemon{}

	for item := range pokemon {
		results = append(results, item)
	}

	return results
}

func (s *Service) listPokemon(offset uint, limit uint) <-chan client.ListPokemonResponse {
	out := make(chan client.ListPokemonResponse)

	go func(offset uint, limit uint) {
		defer close(out)

		for ; ; offset += limit {
			data, err := s.client.ListPokemon(offset, limit)
			if err != nil {
				panic(err.Error())
			}

			out <- data

			if data.Next == "" {
				break
			}
		}
	}(offset, limit)

	return out
}

func (s *Service) getPokemon(listed <-chan client.ListPokemonResponse) <-chan client.DetailPokemonResponse {
	var wg sync.WaitGroup

	out := make(chan client.DetailPokemonResponse)

	go func() {
		defer close(out)

		for data := range listed {
			wg.Add(len(data.Results))

			for _, result := range data.Results {
				go s.getPokemonDetail(&wg, parsePokemonEntryIdFromResponseUrl(result.Url), out)
			}
		}

		wg.Wait()
	}()

	return out
}

func (s *Service) getPokemonDetail(wg *sync.WaitGroup, id uint, out chan<- client.DetailPokemonResponse) {
	defer wg.Done()

	data, err := s.client.GetPokemon(uint(id))
	if err != nil {
		panic(err.Error())
	}

	if data.Order < 0 {
		return
	}

	out <- data
}

func (s *Service) detailPokemonResponseToPokemon(responses <-chan client.DetailPokemonResponse) <-chan Pokemon {
	out := make(chan Pokemon)

	go func() {
		defer close(out)

		for response := range responses {
			out <- mapDetailPokemonResponseToPokemon(response)
		}
	}()

	return out
}

func parsePokemonEntryIdFromResponseUrl(path string) uint {
	comp := strings.Split(path, "/")

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

func mapDetailPokemonResponseToPokemon(response client.DetailPokemonResponse) Pokemon {
	pokemon := Pokemon{
		ID:      response.ID,
		Name:    parsePokemonName(response.Name),
		Number:  response.Order,
		Height:  response.Height,
		Weight:  response.Weight,
		Type:    make([]PokemonType, len(response.Types)),
		Sprite:  response.Sprites.FrontDefault,
		Artwork: response.Sprites.Other.OfficialArtwork.FrontDefault,
	}

	if pokemon.Sprite == "" {
		pokemon.Sprite = response.Sprites.Other.OfficialArtwork.FrontDefault
	}

	for _, stat := range response.Stats {
		switch stat.Stat.Name {
		case "hp":
			pokemon.HP = stat.BaseStat
		case "attack":
			pokemon.Attack = stat.BaseStat
		case "defense":
			pokemon.Defense = stat.BaseStat
		case "special-attack":
			pokemon.SpAttack = stat.BaseStat
		case "special-defense":
			pokemon.SpDefense = stat.BaseStat
		case "speed":
			pokemon.Speed = stat.BaseStat
		}
	}

	for _, respType := range response.Types {
		pokemon.Type[respType.Slot-1] = PokemonType{respType.Type.Name}
	}

	pokemon.BulbapediaURL = fmt.Sprintf("https://bulbapedia.bulbagarden.net/wiki/%v_(Pokémon)", strings.Join(strings.Split(pokemon.Name, " "), "_"))

	return pokemon
}

func parsePokemonName(name string) string {
	output := cases.Title(language.Und).String(strings.Join(strings.Split(name, "-"), " "))

	output = regexp.MustCompile(`^Mr M`).ReplaceAllString(output, "Mr. M")
	output = regexp.MustCompile(`^Mime Jr`).ReplaceAllString(output, "Mime Jr.")
	output = regexp.MustCompile(`^Mr R`).ReplaceAllString(output, "Mr. R")
	output = regexp.MustCompile(`mo O`).ReplaceAllString(output, "mo-o")
	output = regexp.MustCompile(`Porygon Z`).ReplaceAllString(output, "Porygon-Z")
	output = regexp.MustCompile(`Type Null`).ReplaceAllString(output, "Type: Null")
	output = regexp.MustCompile(`Ho Oh`).ReplaceAllString(output, "Ho-Oh")
	output = regexp.MustCompile(`Nidoran F`).ReplaceAllString(output, "Nidoran♀")
	output = regexp.MustCompile(`Nidoran M`).ReplaceAllString(output, "Nidoran♂")
	output = regexp.MustCompile(`Flabebe`).ReplaceAllString(output, "Flabébé")

	return output
}
