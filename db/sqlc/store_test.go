package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	Account1 := createRandomAccount(t)
	Account2 := createRandomAccount(t)
	var amount int64 = 1000

	log.Printf("before >> %d %d", Account1.Balance, Account2.Balance)

	n := 5
	results := make(chan TransferTxResult, n)
	errs := make(chan error, n)

	// wg := &sync.WaitGroup{}
	// wg.Add(n)

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		go func() {
			// defer wg.Done()
			result, err := store.TransferTx(
				context.Background(),
				&TransferTxparams{
					FromAccountID: Account1.ID,
					ToAccountID:   Account2.ID,
					Amount:        amount,
					Description:   nil,
				})
			results <- result
			errs <- err
		}()
	}
	// wg.Wait()
	// close(results)
	// close(errs)

	for i := 0; i < n; i++ {
		err := <-errs
		assert.Nil(t, err)

		result := <-results
		assert.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		assert.Equal(t, Account1.ID, transfer.FromAccountID)
		assert.Equal(t, Account2.ID, transfer.ToAccountID)
		assert.Equal(t, amount, transfer.Amount)
		assert.NotZero(t, transfer.ID)
		assert.NotZero(t, transfer.CreatedAt)

		// check entries
		fromEntry := result.FromEntry
		assert.Equal(t, Account1.ID, fromEntry.AccountID)
		assert.Equal(t, -amount, fromEntry.Amount)
		assert.NotZero(t, fromEntry.ID)
		assert.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		assert.Equal(t, Account2.ID, toEntry.AccountID)
		assert.Equal(t, amount, toEntry.Amount)
		assert.NotZero(t, toEntry.ID)
		assert.NotZero(t, toEntry.CreatedAt)

		// check accounts
		fromAccount := result.FromAccount
		assert.NotEmpty(t, fromAccount)

		toAccount := result.ToAccount
		assert.NotEmpty(t, toAccount)

		log.Printf("tx >> %d %d", fromAccount.Balance, toAccount.Balance)

		// check balances
		diff1 := Account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - Account2.Balance

		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1 > 0)
		assert.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		assert.True(t, k >= 1 && k <= n)
		assert.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), Account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, Account1.Balance-(int64(n)*amount), updatedAccount1.Balance)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), Account2.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedAccount2.Balance, Account2.Balance+(int64(n)*amount))
	log.Printf("after >> %d %d", updatedAccount1.Balance, updatedAccount2.Balance)

}
