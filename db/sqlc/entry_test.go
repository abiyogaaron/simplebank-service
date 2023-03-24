package db

import (
	"context"
	"testing"
	"time"

	"github.com/abiyogaaron/simplebank-service/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(account Account) (Entry, CreateEntryParams, error) {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	return entry, arg, err
}

func TestCreateEntry(t *testing.T) {
	acc, _, _ := createRandomAccount()
	entry, arg, err := createRandomEntry(acc)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	acc, _, _ := createRandomAccount()
	entry, _, _ := createRandomEntry(acc)
	selectedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, selectedEntry)
	require.Equal(t, entry.AccountID, selectedEntry.AccountID)
	require.Equal(t, entry.Amount, selectedEntry.Amount)
	require.WithinDuration(t, entry.CreatedAt, selectedEntry.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	acc, _, _ := createRandomAccount()
	arg := ListEntriesParams{
		Limit:     5,
		Offset:    5,
		AccountID: acc.ID,
	}

	for i := 0; i < 10; i++ {
		createRandomEntry(acc)
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
