package repository

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

// Test_GetFriends used to get friends
func Test_SignUp(t *testing.T) {
	ctx := context.Background()
	repos := NewUserPostgres(db)
	err := repos.SignUp(ctx, &testUserValidData)
	require.NoError(t, err, "get friends error")
}

// Test_GetFriends used to get friends
func Test_SignIn(t *testing.T) {
	ctx := context.Background()
	repos := NewUserPostgres(db)
	err := repos.SignIn(ctx, &testUserValidData)
	require.NoError(t, err, "get friends error")
}
