package models

import (
	"testing"

	"snippet.darieldejesus.com/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	tests := []struct {
		name   string
		userId int
		expect bool
	}{
		{
			name:   "Valid Id",
			userId: 1,
			expect: true,
		},
		{
			name:   "Zero Id",
			userId: 0,
			expect: false,
		},
		{
			name:   "Non-existent Id",
			userId: 7,
			expect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDb(t)
			m := UserModel{db}
			exists, err := m.Exists(test.userId)
			assert.Equal(t, exists, test.expect)
			assert.NilError(t, err)
		})
	}
}
