package migrations

import (
	"strings"
	"testing"
)

func TestGetQueriesFromFile(t *testing.T) {
	pair, err := GetQueriesFromFile("../tests/20221225184300_create_foo.sql")

	if err != nil {
		t.Log("had an error", err)
		t.Fail()
	}

	expectation := strings.Trim(`CREATE TABLE foo (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);`, " \n")

	actual := strings.Trim(pair.Up, " \n")

	if actual != expectation {
		t.Errorf("Expected %v but got %v", expectation, actual)
		t.Fail()
	}
}
