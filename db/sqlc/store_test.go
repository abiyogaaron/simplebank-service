package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	fromAcc, _, _ := createRandomAccount()
	toAcc, _, _ := createRandomAccount()
	fmt.Println(">> Before: ", fromAcc.Balance, toAcc.Balance)

	//run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errMsg := make(chan error)
	resultMsg := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Amount:        amount,
			})

			errMsg <- err
			resultMsg <- result
		}()
	}

	//check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errMsg
		require.NoError(t, err)

		result := <-resultMsg
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAcc.ID, transfer.FromAccountID)
		require.Equal(t, toAcc.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAcc.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, toAcc.ID, ToEntry.AccountID)
		require.Equal(t, amount, ToEntry.Amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), ToEntry.ID)
		require.NoError(t, err)

		//checks account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAcc.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAcc.ID, toAccount.ID)

		//checks account balance
		fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)
		diffFromAcc := fromAcc.Balance - fromAccount.Balance
		diffToAcc := toAccount.Balance - toAcc.Balance
		require.Equal(t, diffFromAcc, diffToAcc)
		require.True(t, diffFromAcc > 0)
		require.True(t, diffFromAcc%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diffFromAcc / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check the final balance
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAcc.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAcc.ID)
	require.NoError(t, err)

	fmt.Println(">> After: ", updatedFromAccount.Balance, updatedToAccount.Balance)
	require.Equal(t, fromAcc.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, toAcc.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	fromAcc, _, _ := createRandomAccount()
	toAcc, _, _ := createRandomAccount()
	fmt.Println(">> Before: ", fromAcc.Balance, toAcc.Balance)

	//run n concurrent transfer transactions
	n := 20
	amount := int64(10)
	errMsg := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := fromAcc.ID
		toAccountId := toAcc.ID

		if i%2 == 1 {
			fromAccountId = toAcc.ID
			toAccountId = fromAcc.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			errMsg <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errMsg
		require.NoError(t, err)
	}

	//check the final balance
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAcc.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAcc.ID)
	require.NoError(t, err)

	fmt.Println(">> After: ", updatedFromAccount.Balance, updatedToAccount.Balance)
	require.Equal(t, fromAcc.Balance, updatedFromAccount.Balance)
	require.Equal(t, toAcc.Balance, updatedToAccount.Balance)
}
