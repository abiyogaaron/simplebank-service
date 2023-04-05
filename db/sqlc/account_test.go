package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/abiyogaaron/simplebank-service/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) (Account, CreateAccountParams, error) {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	return account, arg, err
}

func TestCreateAccount(t *testing.T) {
	account, arg, err := createRandomAccount(t)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	// create account
	account, _, _ := createRandomAccount(t)
	selectedAcc, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, selectedAcc)
	require.Equal(t, account.Owner, selectedAcc.Owner)
	require.Equal(t, account.Balance, selectedAcc.Balance)
	require.Equal(t, account.Currency, selectedAcc.Currency)
	require.WithinDuration(t, account.CreatedAt, selectedAcc.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account, _, _ := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}

	updatedAcc, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)
	require.Equal(t, account.Owner, updatedAcc.Owner)
	require.Equal(t, arg.Balance, updatedAcc.Balance)
	require.Equal(t, account.Currency, updatedAcc.Currency)
	require.WithinDuration(t, account.CreatedAt, updatedAcc.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account, _, _ := createRandomAccount(t)
	errDelete := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, errDelete)

	selectedAcc, errSelect := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, errSelect)
	require.EqualError(t, errSelect, sql.ErrNoRows.Error())
	require.Empty(t, selectedAcc)
}

func TestGetListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
