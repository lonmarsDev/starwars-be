package swservice

import (
	"context"
	"fmt"
	"testing"

	"github.com/lonmarsDev/starwars-be/pkg/database"
	"github.com/lonmarsDev/starwars-be/pkg/model"
	"github.com/lonmarsDev/starwars-be/pkg/swapiclient"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_GetSavedCharacter(t *testing.T) {
	type fields struct {
		Database    *database.DatabaseClient
		SwApiClient *swapiclient.SwApiClient
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Character
		wantErr error
	}{
		{
			name: "happy_path",
			fields: fields{
				Database: &database.DatabaseClient{
					DatabaseClientAction: func() database.DatabaseClientAction {
						ctrl := gomock.NewController(t)

						m := database.NewMockDatabaseClientAction(ctrl)
						mREturn := []model.Character{
							{
								Name:    "star test",
								Film:    []string{"film1", "film2"},
								Vehicle: []string{"v1", "v2"},
							},
						}
						m.EXPECT().GetAllSavedCharacter(context.TODO(), "characters").Return(mREturn, nil)
						return m
					}(),
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: []model.Character{
				{
					Name:    "star test",
					Film:    []string{"film1", "film2"},
					Vehicle: []string{"v1", "v2"},
				},
			},
			wantErr: nil,
		},
		{
			name: "error_path",
			fields: fields{
				Database: &database.DatabaseClient{
					DatabaseClientAction: func() database.DatabaseClientAction {
						ctrl := gomock.NewController(t)

						m := database.NewMockDatabaseClientAction(ctrl)
						mREturn := []model.Character{}
						m.EXPECT().GetAllSavedCharacter(context.TODO(), "characters").Return(mREturn, fmt.Errorf("test error"))
						return m
					}(),
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			want:    []model.Character{},
			wantErr: fmt.Errorf("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Database:    tt.fields.Database,
				SwApiClient: tt.fields.SwApiClient,
			}
			got, err := s.GetSavedCharacter(tt.args.ctx)
			assert.Equal(t, tt.want, got, "should be equal")
			assert.Equal(t, tt.wantErr, err, "should be equal")
		})
	}
}

func TestService_SaveSearch(t *testing.T) {
	type fields struct {
		Database    *database.DatabaseClient
		SwApiClient *swapiclient.SwApiClient
	}
	type args struct {
		ctx   context.Context
		input []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name: "happy_path",
			fields: fields{
				Database: &database.DatabaseClient{
					DatabaseClientAction: func() database.DatabaseClientAction {
						ctrl := gomock.NewController(t)
						m := database.NewMockDatabaseClientAction(ctrl)
						var appInput []interface{}
						appChar := model.Character{
							Name:    "name",
							Film:    []string{"film1", "film2"},
							Vehicle: []string{"v1", "v2"},
						}
						appInput = append(appInput, appChar)
						m.EXPECT().InsertMany(context.TODO(), "characters", appInput).Return("testtest", nil)
						return m
					}(),
				},
			},
			args: args{
				ctx: context.TODO(),
				input: func() []interface{} {
					var appInput []interface{}
					appChar := model.Character{
						Name:    "name",
						Film:    []string{"film1", "film2"},
						Vehicle: []string{"v1", "v2"},
					}
					appInput = append(appInput, appChar)
					return appInput
				}(),
			},
			want:    "testtest",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Database:    tt.fields.Database,
				SwApiClient: tt.fields.SwApiClient,
			}
			got, err := s.SaveSearch(tt.args.ctx, tt.args.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

var testSearchPeople = []byte(`{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "name": "Darth Maul",
            "height": "175",
            "mass": "80",
            "hair_color": "none",
            "skin_color": "red",
            "eye_color": "yellow",
            "birth_year": "54BBY",
            "gender": "male",
            "homeworld": "https://swapi.dev/api/planets/36/",
            "films": [
                "https://swapi.dev/api/films/4/"
            ],
            "species": [
                "https://swapi.dev/api/species/22/"
            ],
            "vehicles": [
                "https://swapi.dev/api/vehicles/42/"
            ],
            "starships": [
                "https://swapi.dev/api/starships/41/"
            ],
            "created": "2014-12-19T18:00:41.929000Z",
            "edited": "2014-12-20T21:17:50.403000Z",
            "url": "https://swapi.dev/api/people/44/"
        }
    ]
}`)

var cineTest = []byte(`{
    "title": "A New Hope",
    "episode_id": 4,
    "opening_crawl": "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden base, have won\r\ntheir first victory against\r\nthe evil Galactic Empire.\r\n\r\nDuring the battle, Rebel\r\nspies managed to steal secret\r\nplans to the Empire's\r\nultimate weapon, the DEATH\r\nSTAR, an armored space\r\nstation with enough power\r\nto destroy an entire planet.\r\n\r\nPursued by the Empire's\r\nsinister agents, Princess\r\nLeia races home aboard her\r\nstarship, custodian of the\r\nstolen plans that can save her\r\npeople and restore\r\nfreedom to the galaxy....",
    "director": "George Lucas",
    "producer": "Gary Kurtz, Rick McCallum",
    "release_date": "1977-05-25",
    "characters": [
        "https://swapi.dev/api/people/1/",
        "https://swapi.dev/api/people/2/",
        "https://swapi.dev/api/people/3/",
        "https://swapi.dev/api/people/4/",
        "https://swapi.dev/api/people/5/",
        "https://swapi.dev/api/people/6/",
        "https://swapi.dev/api/people/7/",
        "https://swapi.dev/api/people/8/",
        "https://swapi.dev/api/people/9/",
        "https://swapi.dev/api/people/10/",
        "https://swapi.dev/api/people/12/",
        "https://swapi.dev/api/people/13/",
        "https://swapi.dev/api/people/14/",
        "https://swapi.dev/api/people/15/",
        "https://swapi.dev/api/people/16/",
        "https://swapi.dev/api/people/18/",
        "https://swapi.dev/api/people/19/",
        "https://swapi.dev/api/people/81/"
    ],
    "planets": [
        "https://swapi.dev/api/planets/1/",
        "https://swapi.dev/api/planets/2/",
        "https://swapi.dev/api/planets/3/"
    ],
    "starships": [
        "https://swapi.dev/api/starships/2/",
        "https://swapi.dev/api/starships/3/",
        "https://swapi.dev/api/starships/5/",
        "https://swapi.dev/api/starships/9/",
        "https://swapi.dev/api/starships/10/",
        "https://swapi.dev/api/starships/11/",
        "https://swapi.dev/api/starships/12/",
        "https://swapi.dev/api/starships/13/"
    ],
    "vehicles": [
        "https://swapi.dev/api/vehicles/4/",
        "https://swapi.dev/api/vehicles/6/",
        "https://swapi.dev/api/vehicles/7/",
        "https://swapi.dev/api/vehicles/8/"
    ],
    "species": [
        "https://swapi.dev/api/species/1/",
        "https://swapi.dev/api/species/2/",
        "https://swapi.dev/api/species/3/",
        "https://swapi.dev/api/species/4/",
        "https://swapi.dev/api/species/5/"
    ],
    "created": "2014-12-10T14:23:31.880000Z",
    "edited": "2014-12-20T19:49:45.256000Z",
    "url": "https://swapi.dev/api/films/1/"
}`)

var testStarShiptest = []byte(`{
    "name": "Sith speeder",
    "model": "FC-20 speeder bike",
    "manufacturer": "Razalon",
    "cost_in_credits": "4000",
    "length": "1.5",
    "max_atmosphering_speed": "180",
    "crew": "1",
    "passengers": "0",
    "cargo_capacity": "2",
    "consumables": "unknown",
    "vehicle_class": "speeder",
    "pilots": [
        "https://swapi.dev/api/people/44/"
    ],
    "films": [
        "https://swapi.dev/api/films/4/"
    ],
    "created": "2014-12-20T10:09:56.095000Z",
    "edited": "2014-12-20T21:30:21.712000Z",
    "url": "https://swapi.dev/api/vehicles/42/"
}`)

func TestService_SearchCharacter(t *testing.T) {
	type fields struct {
		Database    *database.DatabaseClient
		SwApiClient *swapiclient.SwApiClient
	}
	type args struct {
		ctx    context.Context
		filter string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr error
	}{
		{
			name: "happy_path",
			fields: fields{
				SwApiClient: &swapiclient.SwApiClient{
					SwApiClient: func() swapiclient.SwApiClientAction {
						ctrl := gomock.NewController(t)
						m := swapiclient.NewMockSwApiClientAction(ctrl)
						m.EXPECT().SearchPeople(context.TODO(), "Darth Maul").Return(testSearchPeople, nil)
						m.EXPECT().Request(context.TODO(), "https://swapi.dev/api/films/4/").Return(cineTest, nil)
						m.EXPECT().Request(context.TODO(), "https://swapi.dev/api/vehicles/42/").Return(testStarShiptest, nil)
						return m
					}(),
				},
			},
			args: args{
				ctx:    context.TODO(),
				filter: "Darth Maul",
			},
			want: func() interface{} {
				return &model.SearchResult{
					Count: 1,
					Character: []model.Character{
						{
							Name:    "Darth Maul",
							Film:    []string{"A New Hope"},
							Vehicle: []string{"FC-20 speeder bike"},
						},
					},
				}
			}(),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Database:    tt.fields.Database,
				SwApiClient: tt.fields.SwApiClient,
			}
			got, err := s.SearchCharacter(tt.args.ctx, tt.args.filter)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
