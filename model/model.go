package model

import "time"

type FetchSource struct {
	Url  string
	Type []string
}

type Store struct {
	Name      string
	Number    string
	Models    []Model
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	StoreNum  string
	StartTime time.Time
	MessageID int
	ChatID    int
	Enable    bool
	ModelName string
}
