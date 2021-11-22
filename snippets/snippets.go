package snippets

import (
	"database/sql"
	"fmt"
	"sort"
	"sync"

	_ "github.com/lib/pq"
)

type Container interface {
	Insert(Snippet) error
	Find(string) (Snippet, error)
	List() ([]string, error)
	// TODO: Container can return a list of snippets
	// TODO: Container can delete a snippet based off its Name
}

// TODO: Snippet is able to return its in-file text representation
// s.Repr()?
type Snippet struct {
	Name string
	Desc string
	Body string
}

type SnippetsMap struct {
	cntr map[string]Snippet
	sync.Mutex
}

// Create a fresh instance of snippets container.
//
// Allowed types (t): map
func NewSnippetsContainer(t string) (Container, error) {
	/*
		NOTE: At some point, I might want to add case "db" with:
			dns := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
			return NewSnippetsDB("postgres", dns)
	*/
	switch t {
	case "map":
		return NewSnippetsMap(), nil
	default:
		return nil, fmt.Errorf("container type (%s) is not implemented", t)
	}
}

func NewSnippetsMap() *SnippetsMap {
	return &SnippetsMap{
		cntr: make(map[string]Snippet),
	}
}

type SnippetsDB struct {
	db *sql.DB
}

func NewSnippetsDB(dialect string, dsn string) (Container, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &SnippetsDB{db}, nil
}

func (s *SnippetsMap) Insert(snip Snippet) (err error) {
	s.Lock()
	s.cntr[snip.Name], err = snip, nil
	s.Unlock()
	return
}

func (s *SnippetsMap) Find(str string) (Snippet, error) {
	s.Lock()
	defer s.Unlock()
	var snip Snippet
	snip, ok := s.cntr[str]
	if !ok {
		return snip, fmt.Errorf("snippet was not found")
	}
	return snip, nil
}

func (s *SnippetsMap) List() ([]string, error) {
	s.Lock()
	var result []string
	var str string
	for _, v := range s.cntr {
		str = fmt.Sprintf("%s\t%s", v.Name, v.Desc)
		result = append(result, str)
	}
	sort.Strings(result)
	s.Unlock()
	return result, nil
}

func (s *SnippetsDB) Insert(snip Snippet) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()
	query := `
INSERT INTO snippet (
	name,
	"desc",
	body
)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(snip.Name, snip.Desc, snip.Body)
	if i, err := res.RowsAffected(); i == 0 {
		if err != nil {
			return err
		}
		return fmt.Errorf("snippet name already exists in the database")
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *SnippetsDB) Find(str string) (Snippet, error) {
	query := `
SELECT name, "desc", body FROM snippet
where name = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return Snippet{}, err
	}
	row := stmt.QueryRow(str)
	snip := Snippet{}
	err = row.Scan(&snip.Name, &snip.Desc, &snip.Body)
	if err != nil {
		return Snippet{}, err
	}
	return snip, nil
}

func (s *SnippetsDB) List() ([]string, error) {
	var result []string
	query := `SELECT name, "desc" FROM snippet`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return result, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	var str string
	for rows.Next() {
		var name, desc string
		err := rows.Scan(&name, &desc)
		if err != nil {
			return result, err
		}
		str = fmt.Sprintf("%s\t%s", name, desc)
		result = append(result, str)
	}
	return result, nil
}
