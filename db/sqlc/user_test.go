package db

import (
	"context"
	"testing"
	"time"

	"github.com/ilhamgepe/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	password, err := utils.HashPassword(utils.RandomString(12))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username: utils.RandomOwner(),
		Password: password,
		FullName: utils.RandomOwner(),
		Email:    utils.RandomEmail(),
	}
	ctx := context.Background()

	user, err := testQueries.CreateUser(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotEmpty(t, user.Username)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangeAt.IsZero())
	require.WithinDuration(t, time.Now(), user.CreatedAt, time.Minute)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	ctx := context.Background()

	user2, err := testQueries.GetUser(ctx, user1.Username)

	require.NoError(t, err)
	require.Equal(t, user1, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Email, user2.Email)

	require.Equal(t, user1.PasswordChangeAt, user2.PasswordChangeAt)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Minute)
	require.WithinDuration(t, user1.UpdatedAt, user2.UpdatedAt, time.Minute)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullname := utils.RandomOwner()
	newuser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		FullName: utils.NewNullableString(newFullname),
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.Equal(t, oldUser.Username, newuser.Username)
	require.Equal(t, oldUser.Email, newuser.Email)
	require.Equal(t, newFullname, newuser.FullName)
	require.Equal(t, oldUser.Password, newuser.Password)
	require.WithinDuration(t, oldUser.CreatedAt, newuser.CreatedAt, time.Minute)
	require.WithinDuration(t, oldUser.UpdatedAt, newuser.UpdatedAt, time.Minute)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := utils.RandomEmail()
	newuser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Email:    utils.NewNullableString(newEmail),
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.Equal(t, oldUser.Username, newuser.Username)
	require.Equal(t, newEmail, newuser.Email)
	require.Equal(t, oldUser.FullName, newuser.FullName)
	require.Equal(t, oldUser.Password, newuser.Password)
	require.WithinDuration(t, oldUser.CreatedAt, newuser.CreatedAt, time.Minute)
	require.WithinDuration(t, oldUser.UpdatedAt, newuser.UpdatedAt, time.Minute)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword, err := utils.HashPassword(utils.RandomString(12))
	require.NoError(t, err)
	newuser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Password: utils.NewNullableString(newPassword),
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.Equal(t, oldUser.Username, newuser.Username)
	require.Equal(t, oldUser.Email, newuser.Email)
	require.Equal(t, oldUser.FullName, newuser.FullName)
	require.Equal(t, newPassword, newuser.Password)
	require.WithinDuration(t, oldUser.CreatedAt, newuser.CreatedAt, time.Minute)
	require.WithinDuration(t, oldUser.UpdatedAt, newuser.UpdatedAt, time.Minute)
}
