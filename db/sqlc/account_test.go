package db

import (
	"context"
	"testing"
	"time"

	"github.com/ilhamgepe/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  1000,
		Currency: utils.RandomCurrency(),
	}
	ctx := context.Background()

	account, err := testQueries.CreateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.NotZero(t, account.ID)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.WithinDuration(t, time.Now(), account.CreatedAt, time.Minute)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	ctx := context.Background()

	account2, err := testQueries.GetAccount(ctx, account1.ID)

	require.NoError(t, err)
	require.Equal(t, account1, account2)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	ctx := context.Background()

	account2, err := testQueries.UpdateAccount(ctx, arg)

	require.NoError(t, err)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.WithinDuration(t, time.Now(), account2.UpdatedAt, time.Minute)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	ctx := context.Background()

	err := testQueries.DeleteAccount(ctx, account1.ID)

	require.NoError(t, err)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 3; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := GetAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	ctx := context.Background()

	accounts, err := testQueries.GetAccounts(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Equal(t, lastAccount.Owner, accounts[0].Owner)

}
