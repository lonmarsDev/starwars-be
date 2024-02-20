package model

type SearchResult struct {
	Count     int
	Character []Character
}
type Character struct {
	Name    string   `json:"name"`
	Film    []string `json:"films"`
	Vehicle []string `json:"vehicles"`
}
