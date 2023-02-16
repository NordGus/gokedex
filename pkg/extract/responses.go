package extract

// pokemonSpeciesPageResponse
type pokemonSpeciesPageResponse struct {
	Next    string                          `json:"next"`
	Results []pokemonSpeciesPreviewResponse `json:"results"`
}

type pokemonSpeciesPreviewResponse struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// pokemonSpeciesResponse
type pokemonSpeciesResponse struct {
	ID         uint64                             `json:"id"`
	Name       string                             `json:"name"`
	Order      int64                              `json:"order"`
	Generation pokemonSpeciesGenerationResponse   `json:"generation"`
	FlavorText []pokemonSpeciesFlavorTextResponse `json:"flavor_text_entries"`
	Genera     []pokemonSpeciesGenusResponse      `json:"genera"`
	Varieties  []pokemonSpeciesVarietiesResponse  `json:"varieties"`
}

type pokemonSpeciesGenerationResponse struct {
	Name string `json:"name"`
}

type pokemonSpeciesFlavorTextResponse struct {
	Text     string                                   `json:"flavor_text"`
	Language pokemonSpeciesFlavorTextLanguageResponse `json:"language"`
	Version  pokemonSpeciesFlavorTextVersionResponse  `json:"version"`
}

type pokemonSpeciesFlavorTextLanguageResponse struct {
	Name string `json:"name"`
}

type pokemonSpeciesFlavorTextVersionResponse struct {
	Name string `json:"name"`
}

type pokemonSpeciesGenusResponse struct {
	Genus    string                              `json:"genus"`
	Language pokemonSpeciesGenusLanguageResponse `json:"language"`
}

type pokemonSpeciesGenusLanguageResponse struct {
	Name string `json:"name"`
}

type pokemonSpeciesVarietiesResponse struct {
	IsDefault bool                                   `json:"is_default"`
	Pokemon   pokemonSpeciesVarietiesPokemonResponse `json:"pokemon"`
}

type pokemonSpeciesVarietiesPokemonResponse struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// pokemonResponse
type pokemonResponse struct {
	ID      uint64                 `json:"id"`
	Height  uint64                 `json:"height"`
	Weight  uint64                 `json:"weight"`
	Stats   []pokemonStatsResponse `json:"stats"`
	Types   []pokemonTypesResponse `json:"types"`
	Sprites pokemonSpritesResponse `json:"sprites"`
}

type pokemonStatsResponse struct {
	BaseStat uint64              `json:"base_stat"`
	Stat     pokemonStatResponse `json:"stat"`
}

type pokemonStatResponse struct {
	Name string `json:"name"`
}

type pokemonTypesResponse struct {
	Slot uint64              `json:"slot"`
	Type pokemonTypeResponse `json:"type"`
}

type pokemonTypeResponse struct {
	Name string `json:"name"`
}

type pokemonSpritesResponse struct {
	FrontDefault string                      `json:"front_default"`
	Other        pokemonSpritesOtherResponse `json:"other"`
}

type pokemonSpritesOtherResponse struct {
	OfficialArtwork pokemonOfficialArtworkResponse `json:"official-artwork"`
}

type pokemonOfficialArtworkResponse struct {
	FrontDefault string `json:"front_default"`
}
