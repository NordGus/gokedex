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
	workers := 5 * runtime.GOMAXPROCS(0)
	clientWorkers := 5

	pokeapi := http.Client{
		Timeout: time.Millisecond * 3000,
	}

	notion := http.Client{
		Timeout: time.Second * 120,
	}

	pokeapiSem := make(chan bool, clientWorkers)
	notionSem := make(chan bool, clientWorkers)
	globalSem := make(chan bool, workers)

	extractor := extract.NewService(globalSem, &pokeapi, pokeapiSem)
	integrator := integrate.NewService(globalSem, &notion, notionSem, integrate.IntegrationSecret(secret), integrate.DatabaseID(databaseId))

	start := time.Now()
	pokemon := extractor.ExtractPokemon()
	done := integrator.IntegrateToNotion(pokemon)

	<-done

	close(pokeapiSem)
	close(notionSem)
	close(globalSem)

	fmt.Println("Duration:", time.Since(start))
}
