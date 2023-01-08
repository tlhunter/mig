package migrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListFiles(t *testing.T) {
	files, err := ListFiles("../tests/postgres")

	if err != nil {
		t.Log("error listing files", err)
		t.Fail()
		return
	}

	assert.Equal(t, files, []string{
		"20230101120058_add_users_table.sql",
		"20230101120107_add_email_to_users.sql",
	}, "file listing not matching")
}
