package db

import (
	"context"
	"testing"

	"github.com/ilhamgepe/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccount, toAccount Account) Transfer {
	ctx := context.Background()

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(ctx, arg)
	require.NoError(t, err)

	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)
	require.NotZero(t, transfer.ID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount, toAccount)
	// // get account yang di mentransfer
	// currentFromAccount, err := testQueries.GetAccount(ctx, fromAccount.ID)
	// require.NoError(t, err)

	// currentToAccount, err := testQueries.GetAccount(ctx, toAccount.ID)
	// require.NoError(t, err)

	// require.Equal(t, fromAccount.Balance-arg.Amount, currentFromAccount.Balance)
	// require.Equal(t, toAccount.Balance+arg.Amount, currentToAccount.Balance)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	transfer := createRandomTransfer(t, fromAccount, toAccount)

	ctx := context.Background()

	transfer2, err := testQueries.GetTransfer(ctx, transfer.ID)
	require.NoError(t, err)
	require.Equal(t, transfer, transfer2)
}

func TestListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fromAccount, toAccount)
		createRandomTransfer(t, toAccount, fromAccount)
	}

	ctx := context.Background()

	arg := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(ctx, arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == fromAccount.ID || transfer.ToAccountID == toAccount.ID)
	}
}
