package exists_test

import (
	"strings"
	"testing"

	"github.com/augmentable-dev/reqlite/internal/redis/exists"
	_ "github.com/augmentable-dev/reqlite/internal/sqlite"
	"github.com/go-redis/redismock/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.riyazali.net/sqlite"
)

func TestExistsSingle(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	//mockKey := "mykey"
	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		if err := api.CreateFunction("exists", exists.New(rdb)); err != nil {
			return sqlite.SQLITE_ERROR, err
		}
		return sqlite.SQLITE_OK, nil
	})
	mock.ExpectExists("key")
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT EXISTS(yoyo)")

	var s string
	err = row.Scan(&s)
	if err != nil {
		t.Fatal(err, s)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	row = db.QueryRow("SELECT exists('unindexedkey')")
	err = row.Err()
	if err != nil {
		t.Fatal(err)
	}

	err = row.Scan(&s)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(s, "0") != 0 {
		t.Error("returned non zero value when non existent key passed")
	}
}

func TestExistsMultiple(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mockKeys := []string{"mykey1", "mykey2", "mykey3"}
	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		if err := api.CreateFunction("exists", exists.New(rdb)); err != nil {
			return sqlite.SQLITE_ERROR, err
		}
		return sqlite.SQLITE_OK, nil
	})
	mock.ExpectExists(mockKeys...).SetVal(3)
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT exists('mykey1','mykey2','mykey3')")

	var s string
	err = row.Scan(&s)
	if err != nil {
		t.Fatal(err, s)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	row = db.QueryRow("SELECT exists('mykey1','mykey2','notakey1','notakey2')")
	err = row.Err()
	if err != nil {
		t.Fatal(err)
	}

	err = row.Scan(&s)
	if err != nil {
		t.Fatal(err)
	}

	if s != "2" {
		t.Errorf("returned %s value when 2 existent and 2 non-existent real keys passed", s)
	}
}
