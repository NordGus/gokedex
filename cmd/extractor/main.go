package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/NordGus/gokedex/pkg/extract"
	"github.com/NordGus/gokedex/pkg/integrate"
)

func main() {
	secret := os.Getenv("NOTION_INTEGRATION_SECRET")
	databaseId := os.Getenv("NOTION_INTEGRATION_DATABASE_ID")

	pokeapi := http.Client{
		Timeout: time.Millisecond * 3000,
	}

	notion := http.Client{
		Timeout: time.Second * 60,
	}

	extractor := extract.NewService(&pokeapi)
	integrator := integrate.NewService(&notion, integrate.IntegrationSecret(secret), integrate.DatabaseID(databaseId))

	start := time.Now()
	pokemon := extractor.ExtractPokemon()
	finished := integrator.IntegrateToNotion(pokemon)

	<-finished

	fmt.Println("Duration:", time.Since(start))
}
