package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ilhamgepe/simpleBank/utils"
	"github.com/stretchr/testify/assert"
)

func createRandomEntry(t *testing.T, account Account) Entry {

	description := fmt.Sprintf(utils.RandomString(10)+" %v", account.ID)
	var arg = &CreateEntryParams{
		AccountID:   account.ID,
		Amount:      utils.RandomEntriesAmount(),
		Description: &description,
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)

	assert.NoError(t, err)

	assert.Equal(t, arg.AccountID, entry.AccountID)
	assert.Equal(t, arg.Amount, entry.Amount)

	return entry
}
func TestCreateEntry(t *testing.T) {
	createRandomEntry(t, createRandomAccount(t))
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t, createRandomAccount(t))

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)

	assert.NoError(t, err)

	assert.Equal(t, entry.ID, entry2.ID)
	assert.Equal(t, entry.AccountID, entry2.AccountID)
	assert.Equal(t, entry.Amount, entry2.Amount)
	assert.Equal(t, entry.Description, entry2.Description)
	assert.WithinDuration(t, entry.CreatedAt.Time, entry.CreatedAt.Time, time.Second)
}

func TestListEntries(t *testing.T) {
	var entries []Entry
	var account = createRandomAccount(t)
	for i := 0; i < 5; i++ {
		entry := createRandomEntry(t, account)
		entries = append(entries, entry)
	}
	entries2, err := testQueries.ListEntries(context.Background(), &ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	})
	assert.Nil(t, err)
	assert.Len(t, entries2, 5)

	for i, entry := range entries2 {
		assert.NotEmpty(t, entry.ID)
		assert.Equal(t, entry.AccountID, account.ID)
		assert.WithinDuration(t, entries[i].CreatedAt.Time, entry.CreatedAt.Time, time.Second)
	}
}
