package models

import (
	"testing"

	testutils "github.com/sgoldenf/a-place-for-your-thoughts/internal/test_utils"
)

func TestUserModelExists(t *testing.T) {
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := testutils.NewTestDB(t)
			m := UserModel{db}
			exists, err := m.Exists(test.userID)
			testutils.Equal(t, exists, test.want)
			testutils.NilError(t, err)
		})
	}
}
