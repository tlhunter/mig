package migrations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQueriesFromFile(t *testing.T) {
	pair, err := GetQueriesFromFile("../tests/postgres/20230101120058_add_users_table.sql")

	if err != nil {
		t.Log("had an error", err)
		t.Fail()
	}

	expectation := `CREATE TABLE users (
  id serial NOT NULL PRIMARY KEY,
  username varchar(24) UNIQUE
);
INSERT INTO users (username) VALUES ('tlhunter');`

	assert.Equal(t, strings.Trim(pair.Up, " \n"), expectation, "queries aren't equal")
}
