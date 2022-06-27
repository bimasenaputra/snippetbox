package mocks

import (
	"database/sql"
	"time"

	"snippetbox.bimasenaputra/internal/models"
)

var mockSnippet = &models.Snippet{
	ID: 1,
	Title: "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	return 1, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
		case 1:
			return mockSnippet, nil
		default:
			return nil, sql.ErrNoRows
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}

func (m *SnippetModel) GetMaxID() (int, error) {
	return 1, nil
}

func (m *SnippetModel) GetMinID() (int, error) {
	return 1, nil
}

func (m *SnippetModel) NextLatestPaging(id int) ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}

func (m *SnippetModel) PrevLatestPaging(id int) ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}