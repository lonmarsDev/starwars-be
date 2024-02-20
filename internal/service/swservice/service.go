package swservice

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lonmarsDev/starwars-be/pkg/database"
	"github.com/lonmarsDev/starwars-be/pkg/model"
	"github.com/lonmarsDev/starwars-be/pkg/swapiclient"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Database    *database.DatabaseClient
	SwApiClient *swapiclient.SwApiClient
}

func NewService(dbConn string) *Service {
	return &Service{
		Database:    database.NewDatabaseClient(dbConn),
		SwApiClient: swapiclient.NewSwApiClient(),
	}
}

// SearchCharacter is to search starwars characters from swapi
func (s *Service) SearchCharacter(ctx context.Context, filter string) (interface{}, error) {
	charactersRaw, err := s.SwApiClient.SwApiClient.SearchPeople(ctx, filter)
	if err != nil {
		return nil, err
	}

	var charResults swapiclient.CharacterResponse
	err = json.Unmarshal(charactersRaw, &charResults)
	if err != nil {
		fmt.Println("Failed to Unmarshal JSON:", err)
		return nil, err
	}
	var characters []model.Character
	var filmLock sync.Mutex
	var startShipLock sync.Mutex
	for _, character := range charResults.Results {
		var wg sync.WaitGroup

		char := model.Character{Name: character.Name}
		for _, film := range character.Films {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				film, _ := s.getFilmTitle(ctx, id)
				if err != nil {
					logrus.Debugf("failed to get film title: %s", err)
				}
				filmLock.Lock()
				defer filmLock.Unlock()
				char.Film = append(char.Film, film)
			}(film)

		}
		for _, vehicle := range character.Vehicles {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				vehicle, err := s.getVehicleModel(ctx, id)
				if err != nil {
					logrus.Debugf("failed to get vehicle name %s", err)
				}
				startShipLock.Lock()
				defer startShipLock.Unlock()
				char.Vehicle = append(char.Vehicle, vehicle)
			}(vehicle)
		}
		wg.Wait()
		characters = append(characters, char)
	}
	return &model.SearchResult{
		Count:     charResults.Count,
		Character: characters,
	}, nil

}

func (s *Service) getFilmTitle(ctx context.Context, url string) (string, error) {
	var filmObj swapiclient.Film
	filmRaw, err := s.SwApiClient.SwApiClient.Request(ctx, url)
	if err != nil {
		logrus.Errorf("failed to get film title %s", err)
		return "", err
	}
	err = json.Unmarshal(filmRaw, &filmObj)
	if err != nil {
		logrus.Errorf("failed to Unmarshal get film title to JSON %s", err)
		return "", err
	}
	return filmObj.Title, nil
}

func (s *Service) getVehicleModel(ctx context.Context, url string) (string, error) {
	var Obj swapiclient.Vehicle
	ObjRaw, err := s.SwApiClient.SwApiClient.Request(ctx, url)
	if err != nil {
		logrus.Errorf("failed to get sw vehicle model %s", err)
		return "", err
	}
	err = json.Unmarshal(ObjRaw, &Obj)
	if err != nil {
		logrus.Errorf("Failed to Unmarshal get vehicle name to JSON %s", err)
		return "", err
	}
	return Obj.Model, nil
}

func (s *Service) SaveSearch(ctx context.Context, input []interface{}) (string, error) {
	return s.Database.InsertMany(ctx, "characters", input)
}

func (s *Service) GetSavedCharacter(ctx context.Context) ([]model.Character, error) {
	return s.Database.GetAllSavedCharacter(ctx, "characters")
}
