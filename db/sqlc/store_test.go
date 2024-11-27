package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	ctx := context.Background()

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	log.Printf("before transaction: \n fromAccount: %v \n toAccount: %v", fromAccount.Balance, toAccount.Balance)
	// run n concurrent transaction
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	existed := map[int]bool{}
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		// check apakah bener ada data transfernya
		_, err = store.GetTransfer(ctx, transfer.ID)
		require.NoError(t, err)

		// chech FromEntry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(ctx, fromEntry.ID)
		require.NoError(t, err)

		// check ToEntry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(ctx, toEntry.ID)
		require.NoError(t, err)

		// check account
		CurrentFromAccount, err := store.GetAccount(ctx, fromAccount.ID)
		require.NoError(t, err)
		require.Equal(t, CurrentFromAccount.ID, fromAccount.ID)

		CurrentToAccount, err := store.GetAccount(ctx, toAccount.ID)
		require.NoError(t, err)
		require.Equal(t, CurrentToAccount.ID, toAccount.ID)

		log.Printf("in transaction: \n CurrentfromAccount: %v \n CurrenttoAccount: %v", CurrentFromAccount.Balance, CurrentToAccount.Balance)

		// check account Balance
		diff1 := fromAccount.Balance - CurrentFromAccount.Balance //ex: 1000 - 990 = 10 *iteration 1; 1000 - 980 = 20 *iteration 2
		diff2 := CurrentToAccount.Balance - toAccount.Balance     //ex: 1010 - 1000 = 10 *iteration 1; 1020 - 1000 = 20 *iteration 2
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2 > 0)
		require.True(t, diff1%amount == 0) // 10%10 = 0, 20%10 = 0
		require.True(t, diff2%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedFromAccount, err := store.GetAccount(ctx, fromAccount.ID)
	require.NoError(t, err)
	updatedToAccount, err := store.GetAccount(ctx, toAccount.ID)
	require.NoError(t, err)

	log.Printf("after transaction: \n updatedfromAccount: %v \n updatedtoAccount: %v", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, updatedFromAccount.Balance, fromAccount.Balance-(amount*int64(n)))
	require.Equal(t, updatedToAccount.Balance, toAccount.Balance+(amount*int64(n)))

}
