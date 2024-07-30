package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/ilhamgepe/simpleBank/utils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func createRandomAccount(t *testing.T) Account {
	var arg = &CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)

	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}
func TestCreateAccountWithTx(t *testing.T) {
	var arg = &CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}

	tx, err := testDb.BeginTx(context.Background(), pgx.TxOptions{})
	assert.NoError(t, err)
	defer tx.Rollback(context.Background())

	account, err := testQueries.WithTx(tx).CreateAccount(context.Background(), arg)

	assert.NoError(t, err)
	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)

	err = tx.Commit(context.Background())

	assert.NoError(t, err)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestListAccounts(t *testing.T) {
	accounts, err := testQueries.ListAccounts(context.Background(), &ListAccountsParams{
		Limit:  5,
		Offset: 0,
	})

	assert.NoError(t, err)
	assert.Greater(t, len(accounts), 0)
	assert.LessOrEqual(t, len(accounts), 5)

	for _, account := range accounts {
		assert.NotEmpty(t, account.ID)
		assert.NotEmpty(t, account.CreatedAt)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := &UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomBalance(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	assert.NoError(t, err)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, arg.Balance, account2.Balance)
	assert.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	log.Println(account1)
	assert.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, pgx.ErrNoRows.Error())
	assert.Empty(t, account2)
}
