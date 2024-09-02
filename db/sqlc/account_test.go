package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/siavashmirzaeifard/simple_bank/util"
	"github.com/stretchr/testify/require"
)

/*
	we need to connect to database to test create account functionality which is defined in account.go file.
	so the best place to do this connection is in main_test.go file that we have to create it.
	after we complete all those main_test.go configs now we can complete this TestCreateAccount function
*/

/*
	to check the test results, we can use testify, so we need to get it and then add it
	it has several packages, we need to use 'require' package
*/

/*
	we use hard coded args for testing purposes which is not the best practice, instead we can make a util folder
	and generate those data randomly.

	as the next step, we can add test command in Makefile to print verbose and coverage for us, and because we have
	lots of unit tests, then at the end add ./... to run all tests for us
*/

// func TestCreateAccount(t *testing.T) {
// args := CreateAccontParams{
// 	Owner:    "Tom",
// 	Balance:  1000,
// 	Currency: "USD",
// }
// 	// args := CreateAccontParams{
// 	// 	Owner:    "Tom",
// 	// 	Balance:  1000,
// 	// 	Currency: "USD",
// 	// }
// 	args := CreateAccontParams{
// 		Owner:    util.RandomOwner(),
// 		Balance:  int64(util.RandomMoney()),
// 		Currency: util.RandomCurrency(),
// 	}
// 	account, err := testQueries.CreateAccont(context.Background(), args)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, account)
// 	require.Equal(t, args.Owner, account.Owner)
// 	require.Equal(t, args.Balance, account.Balance)
// 	require.Equal(t, args.Currency, account.Currency)
// 	require.NotZero(t, account.ID)
// 	require.NotZero(t, account.CreatedAt)
// }

// ---------------------------------------------------------------------------------------------------------------------

/*
	to test all CRUD functionality we need to create account first, then as we can see on above comment block codes, we
	had it but to make our tests independent from each other, we will make another function to generate this account for
	us and return it to help us use it in all unit tests
*/

func createRandomAccount(t *testing.T) Accounts {
	args := CreateAccontParams{
		Owner:    util.RandomOwner(),
		Balance:  int64(util.RandomMoney()),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccont(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      account1.ID,
		Balance: int64(util.RandomMoney()),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, args.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	// below code means skip the first 5 records and return last 5 records
	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
