package models

import (
	"database/sql"
	"errors"
	"time"
	"snippetbox.bimasenaputra/internal/util"
)

type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModelInterface interface {
	Insert(string, string, int) (int, error)
	Get(int) (*Snippet, error)
	Latest() ([]*Snippet, error)
	GetMaxID() (int, error)
	GetMinID() (int, error)
	NextLatestPaging(int) ([]*Snippet, error)
	PrevLatestPaging(int) ([]*Snippet, error)
	LatestContainsTitle(string) ([]*Snippet, error)
	GetMaxIDByTitle(string) (int, error)
	GetMinIDByTitle(string) (int, error)
	NextLatestContainsTitle(int, string) ([]*Snippet, error)
	PrevLatestContainsTitle(int, string) ([]*Snippet, error)  
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {

	stmt := `INSERT INTO SNIPPETS (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {

	stmt := `SELECT * FROM SNIPPETS 
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	s := &Snippet{}

	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {

	stmt := `SELECT * FROM SNIPPETS 
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) GetMaxID() (int, error) {
	stmt := `SELECT MAX(id) FROM SNIPPETS
	WHERE expires > UTC_TIMESTAMP()`

	var id *int

	err := m.DB.QueryRow(stmt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return *id, nil
}

func (m *SnippetModel) GetMinID() (int, error) {
	stmt := `SELECT MIN(id) FROM SNIPPETS
	WHERE expires > UTC_TIMESTAMP()`

	var id *int

	err := m.DB.QueryRow(stmt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return *id, nil
}

func (m *SnippetModel) NextLatestPaging(id int) ([]*Snippet, error) {
	stmt := `SELECT * FROM SNIPPETS
	WHERE id < ? AND expires > UTC_TIMESTAMP()
	ORDER BY id DESC LIMIT 10`

	result, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	snippets := []*Snippet{}

	for result.Next() {
		s := &Snippet{}

		err := result.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) PrevLatestPaging(id int) ([]*Snippet, error) {
	stmt := `SELECT * FROM SNIPPETS
	WHERE id > ? AND expires > UTC_TIMESTAMP()
	ORDER BY id LIMIT 10`

	result, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	snippets := []*Snippet{}

	for result.Next() {
		s := &Snippet{}

		err := result.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	util.Reverse(snippets)

	return snippets, nil
}

func (m *SnippetModel) LatestContainsTitle(title string) ([]*Snippet, error) {
	stmt := `SELECT * FROM SNIPPETS
	WHERE expires > UTC_TIMESTAMP() AND MATCH(title) AGAINST(?)
	ORDER BY id DESC LIMIT 10`

	result, err := m.DB.Query(stmt, title)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	snippets := []*Snippet{}

	for result.Next() {
		s := &Snippet{}

		err := result.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) GetMaxIDByTitle(title string) (int, error) {
	stmt := `SELECT MAX(id) FROM SNIPPETS
	WHERE expires > UTC_TIMESTAMP() AND MATCH(title) AGAINST(?)`

	var id *int

	err := m.DB.QueryRow(stmt, title).Scan(&id)
	if err != nil {
		return 0, err
	} else if id == nil {
		return 0, ErrNoRecord
	}

	return *id, err
}

func (m *SnippetModel) GetMinIDByTitle(title string) (int, error) {
	stmt := `SELECT MIN(id) FROM SNIPPETS
	WHERE expires > UTC_TIMESTAMP() AND MATCH(title) AGAINST(?)`

	var id *int

	err := m.DB.QueryRow(stmt, title).Scan(&id)
	if err != nil {
		return 0, err
	} else if id == nil {
		return 0, ErrNoRecord
	}

	return *id, err
}

func (m *SnippetModel) NextLatestContainsTitle(id int, title string) ([]*Snippet, error) {
	stmt := `SELECT * FROM SNIPPETS
	WHERE id < ? AND expires > UTC_TIMESTAMP() AND MATCH(title) AGAINST(?)
	ORDER BY id DESC LIMIT 10`

	result, err := m.DB.Query(stmt, id, title)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	snippets := []*Snippet{}

	for result.Next() {
		s := &Snippet{}

		err := result.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) PrevLatestContainsTitle(id int, title string) ([]*Snippet, error) {
	stmt := `SELECT * FROM SNIPPETS
	WHERE id > ? AND expires > UTC_TIMESTAMP() AND MATCH(title) AGAINST(?)
	ORDER BY id LIMIT 10`

	result, err := m.DB.Query(stmt, id, title)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	snippets := []*Snippet{}

	for result.Next() {
		s := &Snippet{}

		err := result.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	util.Reverse(snippets)

	return snippets, nil
}