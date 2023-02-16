package integrate

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/NordGus/gokedex/pkg/extract"
)

type Service struct {
	client client
}

func NewService(conn Connection, secret IntegrationSecret, databaseId DatabaseID) Service {
	return Service{
		client: newClient(conn, secret, databaseId),
	}
}

func (s *Service) IntegrateToNotion(in <-chan extract.Pokemon) <-chan struct{} {
	out := make(chan struct{})

	go func(in <-chan extract.Pokemon, out chan<- struct{}) {
		defer close(out)
		pokemon := []extract.Pokemon{}

		for poke := range in {
			pokemon = append(pokemon, poke)
		}

		random := rand.New(rand.NewSource(time.Now().Unix()))
		sample := pokemon[random.Intn(len(pokemon))]

		fmt.Println("Processed Pokémon:", len(pokemon))
		fmt.Println("Sample Pokémon:")
		fmt.Println("\tID:", sample.ID)
		fmt.Println("\tName:", sample.Name)
		fmt.Println("\tNumber:", sample.Number)
		fmt.Printf("\tType: %+v\n", sample.Type)
		fmt.Println("\tHeight:", sample.Height)
		fmt.Println("\tWeight:", sample.Weight)
		fmt.Println("\tHP:", sample.HP)
		fmt.Println("\tAttack:", sample.Attack)
		fmt.Println("\tDefense:", sample.Defense)
		fmt.Println("\tSpAttack:", sample.SpAttack)
		fmt.Println("\tSpDefense:", sample.SpDefense)
		fmt.Println("\tSpeed:", sample.Speed)
		fmt.Printf("\tCategory: %+v\n", sample.Category)
		fmt.Println("\tGeneration:", sample.Generation)
		fmt.Println("\tSprite:", sample.Sprite)
		fmt.Println("\tArtwork:", sample.Artwork)
		fmt.Println("\tBulbapediaURL:", sample.BulbapediaURL)
		fmt.Println("\tFlavorText:")
		for _, text := range sample.FlavorText {
			fmt.Println("\t\tVersion:", text.Version)
			fmt.Println("\t\tText:", text.Text)
		}
	}(in, out)

	return out
}
