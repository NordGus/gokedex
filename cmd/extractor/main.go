package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/NordGus/gokedex/pkg/extract"
	"github.com/NordGus/gokedex/pkg/integrate"
)

func main() {
	secret := os.Getenv("NOTION_INTEGRATION_SECRET")
	databaseId := os.Getenv("NOTION_INTEGRATION_DATABASE_ID")
	workers := 3 * runtime.GOMAXPROCS(0)

	pokeapi := http.Client{
		Timeout: time.Millisecond * 3000,
	}

	notion := http.Client{
		Timeout: time.Second * 120,
	}

	extractor := extract.NewService(&pokeapi)
	integrator := integrate.NewService(&notion, integrate.IntegrationSecret(secret), integrate.DatabaseID(databaseId))
	limits := make(chan bool, workers)

	start := time.Now()
	pokemon := extractor.ExtractPokemon(limits)
	done := integrator.IntegrateToNotion(pokemon, limits)

	<-done

	close(limits)

	fmt.Println("Duration:", time.Since(start))
}
