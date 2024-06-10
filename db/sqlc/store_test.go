package db

import (
	"context"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	Account1 := createRandomAccount(t)
	Account2 := createRandomAccount(t)

	n := 1
	results := make(chan TransferTxResult, n)
	errs := make(chan error, n)

	wg := &sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			result, err := store.TransferTx(context.Background(), &TransferTxparams{
				FromAccountID: Account1.ID,
				ToAccountID:   Account2.ID,
				Amount:        10,
				Description:   nil,
			})
			results <- result
			errs <- err
		}()
	}
	wg.Wait()
	close(results)
	close(errs)

	for err := range errs {
		log.Println(err)
		assert.Nil(t, err)
	}

	for result := range results {
		log.Println("result ", result.ToEntry.ID, result.ToEntry.AccountID, result.ToAccount.ID)
		assert.NotEmpty(t, result)
	}
}
