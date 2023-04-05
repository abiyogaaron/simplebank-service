package db

import (
	"context"
	"testing"
	"time"

	"github.com/abiyogaaron/simplebank-service/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	selectedUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, selectedUser)

	require.Equal(t, user.Username, selectedUser.Username)
	require.Equal(t, user.HashedPassword, selectedUser.HashedPassword)
	require.Equal(t, user.FullName, selectedUser.FullName)
	require.Equal(t, user.Email, selectedUser.Email)

	require.WithinDuration(t, user.PasswordChangedAt, selectedUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user.CreatedAt, selectedUser.CreatedAt, time.Second)
}
