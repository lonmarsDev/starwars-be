package swapiclient

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/http2"
)

const apiBaseURL string = "https://swapi.dev/api"

// for more info, please see documentation page at
// https://swapi.dev/documentation

//go:generate mockgen -destination mock_client.go -package swapiclient github.com/lonmarsDev/starwars-be/pkg/swapiclient  SwApiClientAction
type SwApiClientAction interface {
	SearchPeople(ctx context.Context, name string) (resp []byte, err error)
	Request(ctx context.Context, url string) (resp []byte, err error)
}

type SwApiClient struct {
	SwApiClient SwApiClientAction
}

func NewSwApiClient() *SwApiClient {
	client := &http.Client{
		Transport: &http2.Transport{},
		Timeout:   20 * time.Second,
	}
	return &SwApiClient{
		SwApiClient: &SwApiClientImp{
			BaseUrl: apiBaseURL,
			Client:  client,
		},
	}
}

type SwApiClientImp struct {
	BaseUrl string
	Client  *http.Client
}

func (r *SwApiClientImp) SearchPeople(ctx context.Context, name string) ([]byte, error) {
	params := url.Values{}
	params.Add("search", name)
	bb, err := url.JoinPath(r.BaseUrl, "people", "")
	if err != nil {
		return nil, err
	}
	apiUrl := bb + "?" + params.Encode()
	if err != nil {
		return nil, err
	}
	return r.swapiBodyReader(ctx, apiUrl)
}

func (r *SwApiClientImp) Request(ctx context.Context, url string) ([]byte, error) {
	return r.swapiBodyReader(ctx, url)
}

func (r *SwApiClientImp) swapiBodyReader(ctx context.Context, swapiUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", swapiUrl, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := r.Client.Do(req)
	defer resp.Body.Close()
	// Read the response body
	return io.ReadAll(resp.Body)
}
