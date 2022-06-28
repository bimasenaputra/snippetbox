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
	switch id {
	case 2:
		return []*models.Snippet{mockSnippet}, nil
	default:
		return nil, nil
	}
}

func (m *SnippetModel) PrevLatestPaging(id int) ([]*models.Snippet, error) {
	switch id {
	case 0:
		return []*models.Snippet{mockSnippet}, nil
	default:
		return nil, nil
	}
}

func (m *SnippetModel) LatestContainsTitle(title string) ([]*models.Snippet, error) {
	switch title {
	case "Old":
		return []*models.Snippet{mockSnippet}, nil
	default:
		return nil, nil
	}
}

func (m *SnippetModel) GetMaxIDByTitle(title string) (int, error) {
	switch title {
	case "Old":
		return 1, nil
	default:
		return 0, nil
	}
}

func (m *SnippetModel) GetMinIDByTitle(title string) (int, error) {
	switch title {
	case "Old":
		return 1, nil
	default:
		return 0, nil
	}
}

func (m *SnippetModel) NextLatestContainsTitle(id int, title string) ([]*models.Snippet, error) {
	switch id {
	case 2:
		switch title {
		case "Old":
			return []*models.Snippet{mockSnippet}, nil
		default:
			return nil, nil
		}
	default:
		return nil, nil
	}
}

func (m *SnippetModel) PrevLatestContainsTitle(id int, title string) ([]*models.Snippet, error) {
	switch id {
	case 0:
		switch title {
		case "Old":
			return []*models.Snippet{mockSnippet}, nil
		default:
			return nil, nil
		}
	default:
		return nil, nil
	}
}