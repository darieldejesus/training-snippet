package mocks

import (
	"time"

	"snippet.darieldejesus.com/internal/models"
)

var mockSnippet = &models.Snippet{
	ID:      7,
	Title:   "Lorem ipsum",
	Content: "Lorem ipsum dolor sit amet...",
	Created: time.Date(2022, 12, 7, 11, 12, 0, 0, time.UTC),
	Expires: time.Date(2023, 12, 7, 11, 12, 0, 0, time.UTC),
}

type SnippetModel struct {
	Err error
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	return 8, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	switch id {
	case 7:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return []*models.Snippet{mockSnippet}, nil
}
