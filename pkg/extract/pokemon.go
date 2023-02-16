package extract

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
	Category      []string
	Generation    string
	Sprite        string
	Artwork       string
	BulbapediaURL string
	FlavorText    []PokemonFlavorText
}

type PokemonType struct {
	Name string
}

type PokemonFlavorText struct {
	Text    string
	Version string
}

func mapPokemon(species pokemonSpeciesResponse, details pokemonResponse) Pokemon {
	output := Pokemon{
		ID:            details.ID,
		Name:          parsePokemonName(species.Name),
		Number:        species.Order,
		Height:        details.Height,
		Weight:        details.Weight,
		Type:          parsePokemonType(details.Types),
		Category:      parsePokemonCategory(species.Genera),
		Generation:    parsePokemonGeneration(species.Generation.Name),
		Sprite:        parsePokemonSprite(details.Sprites),
		Artwork:       details.Sprites.Other.OfficialArtwork.FrontDefault,
		BulbapediaURL: parsePokemonBulbapediaURL(parsePokemonName(species.Name)),
		FlavorText:    parsePokemonFlavorText(species.FlavorText),
	}

	for _, stat := range details.Stats {
		switch stat.Stat.Name {
		case "hp":
			output.HP = stat.BaseStat
		case "attack":
			output.Attack = stat.BaseStat
		case "defense":
			output.Defense = stat.BaseStat
		case "special-attack":
			output.SpAttack = stat.BaseStat
		case "special-defense":
			output.SpDefense = stat.BaseStat
		case "speed":
			output.Speed = stat.BaseStat
		}
	}

	return output
}

func parsePokemonName(raw string) string {
	output := cases.Title(language.Und).String(regexp.MustCompile(`-`).ReplaceAllString(raw, " "))
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

func parsePokemonType(types []pokemonTypesResponse) []PokemonType {
	output := make([]PokemonType, len(types))

	for _, t := range types {
		output[t.Slot-1] = PokemonType{t.Type.Name}
	}

	return output
}

func parsePokemonCategory(genera []pokemonSpeciesGenusResponse) []string {
	output := []string{}

	for _, genus := range genera {
		if genus.Language.Name == "en" {
			output = append(output, genus.Genus)
		}
	}

	return output
}

func parsePokemonGeneration(raw string) string {
	splited := strings.Split(raw, "-")
	output := cases.Upper(language.Und).String(splited[len(splited)-1])

	return output
}

func parsePokemonFlavorText(flavors []pokemonSpeciesFlavorTextResponse) []PokemonFlavorText {
	output := []PokemonFlavorText{}

	for _, flavor := range flavors {
		if flavor.Language.Name == "en" {
			output = append(output, PokemonFlavorText{
				Text:    parsePokemonFlavorTextEntry(flavor.Text),
				Version: parsePokemonFlavorTextVersionName(flavor.Version.Name),
			})
		}
	}

	return output
}

func parsePokemonSprite(sprite pokemonSpritesResponse) string {
	if sprite.FrontDefault == "" {
		return sprite.Other.OfficialArtwork.FrontDefault
	}

	return sprite.FrontDefault
}

func parsePokemonBulbapediaURL(raw string) string {
	output := regexp.MustCompile(`\s`).ReplaceAllString(raw, "_")
	return fmt.Sprintf("https://bulbapedia.bulbagarden.net/wiki/%v_(Pokémon)", output)
}

func parsePokemonFlavorTextEntry(raw string) string {
	return regexp.MustCompile(`\s`).ReplaceAllString(raw, " ")
}

func parsePokemonFlavorTextVersionName(raw string) string {
	return cases.Title(language.Und).String(strings.Join(strings.Split(raw, "-"), " "))
}
