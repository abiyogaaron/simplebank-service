package db

import (
	"context"
	"testing"
	"time"

	"github.com/abiyogaaron/simplebank-service/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(accFrom Account, accTo Account) (Transfer, CreateTransferParams, error) {
	arg := CreateTransferParams{
		FromAccountID: accFrom.ID,
		ToAccountID:   accTo.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	return transfer, arg, err
}

func TestCreateTransfer(t *testing.T) {
	accFrom, _, _ := createRandomAccount()
	accTo, _, _ := createRandomAccount()

	transfer, arg, err := createRandomTransfer(accFrom, accTo)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransfer(t *testing.T) {
	accFrom, _, _ := createRandomAccount()
	accTo, _, _ := createRandomAccount()
	transfer, _, _ := createRandomTransfer(accFrom, accTo)

	selectedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, selectedTransfer)
	require.Equal(t, transfer.FromAccountID, selectedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, selectedTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, selectedTransfer.Amount)
	require.WithinDuration(t, transfer.CreatedAt, selectedTransfer.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	accFrom, _, _ := createRandomAccount()
	accTo, _, _ := createRandomAccount()

	arg := ListTransfersParams{
		FromAccountID: accFrom.ID,
		ToAccountID:   accTo.ID,
		Limit:         5,
		Offset:        5,
	}

	for i := 0; i < 10; i++ {
		createRandomTransfer(accFrom, accTo)
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.Equal(t, accFrom.ID, transfer.FromAccountID)
		require.Equal(t, accTo.ID, transfer.ToAccountID)
	}
}
