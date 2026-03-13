package pokemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type pokeAPIResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Types  []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

type pokeAPISpeciesResponse struct {
	Color struct {
		Name string `json:"name"`
	} `json:"color"`
	Habitat struct {
		Name string `json:"name"`
	} `json:"habitat"`
	Generation struct {
		Name string `json:"name"`
	} `json:"generation"`
	EvolvesFromSpecies *struct {
		Name string `json:"name"`
	} `json:"evolves_from_species"`
}

func (c *Client) GetPokemon(ctx context.Context, nameOrID string) (*Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, nameOrID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("pokemon %q not found", nameOrID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var apiResp pokeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	species, err := c.GetSpecies(ctx, apiResp.ID)
	if err != nil {
		return nil, fmt.Errorf("get species: %w", err)
	}

	pokemon := &Pokemon{
		ID:        apiResp.ID,
		Name:      apiResp.Name,
		Height:    apiResp.Height,
		SpriteURL: apiResp.Sprites.FrontDefault,
		Color:     species.Color.Name,
		Habitat:   species.Habitat.Name,
	}

	for _, t := range apiResp.Types {
		pokemon.Types = append(pokemon.Types, PokemonType{Name: t.Type.Name})
	}

	pokemon.Generation = parseGeneration(species.Generation.Name)

	return pokemon, nil
}

func (c *Client) GetSpecies(ctx context.Context, id int) (*pokeAPISpeciesResponse, error) {
	url := fmt.Sprintf("%s/pokemon-species/%d", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var speciesResp pokeAPISpeciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&speciesResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &speciesResp, nil
}

func parseGeneration(name string) int {
	generations := map[string]int{
		"generation-i":    1,
		"generation-ii":   2,
		"generation-iii":  3,
		"generation-iv":   4,
		"generation-v":    5,
		"generation-vi":   6,
		"generation-vii":  7,
		"generation-viii": 8,
		"generation-ix":   9,
	}

	if gen, ok := generations[name]; ok {
		return gen
	}
	return 0
}
