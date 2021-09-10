package snippets

import (
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestContainerFindMethod(t *testing.T) {
	var ss Container = make(SnippetsMap)
	ss.Insert(Snippet{
		Name: "anonfunc",
		Desc: "anonymous function in the Go programming language",
		Body: "func () {${1:body}}()",
	})
	_, err := ss.Find("anonfunc")
	if err != nil {
		t.Errorf("snippets fails to recover existing snippet")
	}
}

func TestSnippetsMapInsert(t *testing.T) {
	ss := make(SnippetsMap)
	err := ss.Insert(Snippet{"name", "desc", "body"})
	if err != nil {
		t.Error("Insert() fails to insert Snippet to map")
	}
}

func TestSnippetsMapFind(t *testing.T) {
	ss := make(SnippetsMap)
	ss["func"] = Snippet{"func", "Go function", "func ${1:name} () {}"}
	_, err := ss.Find("func")
	if err != nil {
		t.Error("existing snippet signature could not be retrieved")
	}
}

func TestSnippetsMapList(t *testing.T) {
	ss := SnippetsMap{
		"func":   {"func", "Go function", "func() {}"},
		"struct": {"struct", "Go struct", "type struct {}"},
		"map":    {"map", "Go map", "map[string]string"},
	}
	want := []string{"func\tGo function", "map\tGo map", "struct\tGo struct"}
	if has, err := ss.List(); !reflect.DeepEqual(has, want) || err != nil {
		t.Errorf("want: %v; has: %v", want, has)
	}
}

func TestMockDBInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error occurred when mocking a database: %s", err)
	}
	defer db.Close()

	stmt := `INSERT INTO snippet`
	args := []driver.Value{"func", "simple func", "def main(): pass"}
	result := sqlmock.NewResult(0, 1)
	mock.ExpectBegin()
	mock.ExpectPrepare(stmt)
	mock.ExpectExec(stmt).WithArgs(args...).WillReturnResult(result)
	mock.ExpectCommit()

	s := SnippetsDB{
		db: db,
	}
	snip := Snippet{Name: "func", Desc: "simple func", Body: "def main(): pass"}
	if err := s.Insert(snip); err != nil {
		t.Errorf("unexpected error when executing INSERT: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were left unfulfilled: %s", err)
	}
}

func TestDBFind(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error occurred when mocking a database: %s", err)
	}
	defer db.Close()

	stmt := `SELECT (.*) FROM snippet`
	rows := sqlmock.NewRows([]string{"name", "desc", "body"}).
		AddRow("func", "simple func", "def main(): pass")

	mock.ExpectPrepare(stmt).
		ExpectQuery().
		WillReturnRows(rows)

	s := SnippetsDB{
		db: db,
	}
	if _, err := s.Find("func"); err != nil {
		t.Errorf("unexpected error when executing SELECT: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were left unfulfilled: %s", err)
	}
}

func TestDBList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error occurred when mocking a database: %s", err)
	}
	defer db.Close()

	stmt := `SELECT (.*) FROM snippet`
	rows := sqlmock.NewRows([]string{"name", "desc"}).
		AddRow("func", "simple func").
		AddRow("gfunc", "function in Go")

	mock.ExpectPrepare(stmt).
		ExpectQuery().
		WillReturnRows(rows)

	s := SnippetsDB{
		db: db,
	}
	if _, err := s.List(); err != nil {
		t.Errorf("unexpected error when executing SELECT: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were left unfulfilled: %s", err)
	}
}
