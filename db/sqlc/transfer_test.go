package db

import (
	"context"
	"testing"
	"time"

	"github.com/ilhamgepe/simpleBank/utils"
	"github.com/stretchr/testify/assert"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := &CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomBalance(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	assert.Nil(t, err)

	assert.Equal(t, arg.Amount, transfer.Amount)
	assert.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	assert.Equal(t, arg.ToAccountID, transfer.ToAccountID)

	return transfer
}
func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	tf := createRandomTransfer(t, account1, account2)

	tf2, err := testQueries.GetTransfer(context.Background(), tf.ID)

	assert.Nil(t, err)

	assert.Equal(t, tf.ID, tf2.ID)
	assert.Equal(t, tf.FromAccountID, tf2.FromAccountID)
	assert.Equal(t, tf.ToAccountID, tf2.ToAccountID)
	assert.Equal(t, tf.Amount, tf2.Amount)
	assert.WithinDuration(t, tf.CreatedAt.Time, tf2.CreatedAt.Time, time.Second)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	var transfers []Transfer
	for i := 0; i < 5; i++ {
		tf := createRandomTransfer(t, account1, account2)
		transfers = append(transfers, tf)
	}
	arg := &ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        0,
	}
	transfers2, err := testQueries.ListTransfers(context.Background(), arg)

	assert.Nil(t, err)

	assert.Len(t, transfers, 5)
	assert.Len(t, transfers2, 5)

	for i, transfer := range transfers2 {
		assert.NotEmpty(t, transfer.ID)
		assert.Equal(t, transfers[i].ID, transfer.ID)
		assert.Equal(t, transfers[i].FromAccountID, transfer.FromAccountID)
		assert.Equal(t, transfers[i].ToAccountID, transfer.ToAccountID)
		assert.Equal(t, transfers[i].Amount, transfer.Amount)
		assert.WithinDuration(t, transfers[i].CreatedAt.Time, transfer.CreatedAt.Time, time.Second)
	}
}
