package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

//		func TestTransferTx(t *testing.T) {
//			// as NewStore function get a new sql.DB object, we need to change our main_test database and add a new global varianle
//			//		for it, then pass it into new store
//			store := NewStore(testDB)
//			account1 := createRandomAccount(t)
//			account2 := createRandomAccount(t)
//			fmt.Println(">> before: ", account1.Balance, account2.Balance)
//			// the best practice to handle a database transaction is to run it in several concurrent Go Routines
//			// run n concurrent transfer transactions
//			n := 2
//			amount := int64(10)
//			// it is a wrong way and if we want to check our result and error inside on the below function and
//			//		we can not just use testify require functions to check results and errors right there in go routines,
//			//		because that function is running inside a different Go Routnine from the one that our TestTransferTx
//			// 		function is running on, so there is no guarantee that it will stop the whole test if a condition is not
//			//  	satisfied
//			// so the correct way to check the result and err is to send them back to the main Goroutine that our test is running on
//			//		and send them back there, and that is why we define channels to get connect concurrent Goroutines
//			errs := make(chan error)
//			results := make(chan TransferTxResult)
//			for i := 0; i < n; i++ {
//				go func() {
//					result, err := store.TransferTx(context.Background(), TransferTxParams{
//						FromAccountID: account1.ID,
//						ToAccountID:   account2.ID,
//						Amount:        amount,
//					})
//					errs <- err
//					results <- result
//				}()
//			}
//			// we made n transfers in upper lines, now we want to check result
//			// check result
//			for i := 0; i < n; i++ {
//				err := <-errs
//				require.NoError(t, err)
//				result := <-results
//				require.NotEmpty(t, result)
//				// explained why we need this variable in comments below on line 111
//				existed := make(map[int]bool)
//				// check transfer
//				transfer := result.Transfers
//				require.NotEmpty(t, transfer)
//				require.Equal(t, account1.ID, transfer.FromAccountID)
//				require.Equal(t, account2.ID, transfer.ToAccountID)
//				require.Equal(t, amount, transfer.Amount)
//				require.NotZero(t, transfer.ID)
//				require.NotZero(t, transfer.CreatedAt)
//				_, err = store.GetTransfer(context.Background(), transfer.ID)
//				require.NoError(t, err)
//				// check entries
//				fromEntry := result.FromEntry
//				require.NotEmpty(t, fromEntry)
//				require.Equal(t, account1.ID, fromEntry.AccountID)
//				require.Equal(t, -amount, fromEntry.Amount)
//				require.NotZero(t, fromEntry.ID)
//				require.NotZero(t, fromEntry.CreatedAt)
//				_, err = store.GetEntry(context.Background(), fromEntry.ID)
//				require.NoError(t, err)
//				toEntry := result.ToEntry
//				require.NotEmpty(t, toEntry)
//				require.Equal(t, account2.ID, toEntry.AccountID)
//				require.Equal(t, amount, toEntry.Amount)
//				require.NotZero(t, toEntry.ID)
//				require.NotZero(t, toEntry.CreatedAt)
//				_, err = store.GetEntry(context.Background(), toEntry.ID)
//				require.NoError(t, err)
//				/*
//					till now, we first wrote codes and then unit tests, since now we will do it in TDD
//					 so first we wrote our tests with logics in store_test.go, then will come back here
//					 to implement logics that pass those tests. so first we check accounts and diff1 and diff2
//					 to check they are fine, then after the for loop we should check updated balance for both accounts
//				*/
//				// check  accounts
//				fromAccount := result.FromAccount
//				require.NotZero(t, fromAccount)
//				require.Equal(t, fromAccount.ID, account1.ID)
//				toAccount := result.ToAccount
//				require.NotEmpty(t, toAccount)
//				require.Equal(t, toAccount.ID, account2.ID)
//				fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)
//				diff1 := account1.Balance - fromAccount.Balance
//				diff2 := toAccount.Balance - account2.Balance
//				require.Equal(t, diff1, diff2)
//				require.True(t, diff1 > 0)
//				require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount
//				// k must be unique for each transaction, 1 for the first transaction, 2 for the second, ...
//				//		and n for the nth transaction. in order to check this, we need to declare a new variable
//				//		called existed of type map[int]bool
//				k := int(diff1 / amount)
//				require.True(t, k >= 0 && k <= n)
//				require.NotContains(t, existed, k)
//				existed[k] = true
//			}
//			updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
//			require.NoError(t, err)
//			updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
//			require.NoError(t, err)
//			require.Equal(t, updatedAccount1.Balance, account1.Balance-int64(n)*amount)
//			require.Equal(t, updatedAccount2.Balance, account2.Balance+int64(n)*amount)
//			fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)
//			/*
//				here is our code's first version, if we run the test it will fail due to deadlock happen in our implementation.
//				 	if we go more deep on deadlocks, and search for deadlocks in postgres database we will find a lot of queries which
//					we can run and find our deadlocks, so after debugging we will find the deadlock is because of the relationship
//					between Accounts table and Transfers table plus Entries table due to foreign keys in migrate up for transfers
//					table and entities table. to fix that a simple way which is not efficient is remove those ALTER TABLE statements
//					and migrate down and up once again.
//					but before to find where is the deadlock from, we can add more debuggers to find which parts of queries have deadlock.
//					for that purpose we have to assign a name for each transaction and pass it into the TransferTx function via the
//					context arguments, to do that inside for loop before the TransferTx function we can crreate a txName variable like this:
//				txName:=fmt.Sprintf("tx %d",i+1)
//					then inside the Goroutines, instead of passing the context.Background(), we should make a context variable in top of the
//					TransferTx function, and pass it into the TransferTx with that transaction name, but as we can see in context.WithValue
//					docs, we find that the second arg which is the provided context key should not be one those built-in types, therefore we should
//					make an empty struct for it in store.go file below the other structs like this: var txKey := struct{}{}
//						the second brackets means that we want to make a new empty object of type that.
//					then define our new context and pass those txKey and txName to the function like this:
//				ctx := context.WithValue(context.Background(), txKey, txName)
//					after this step context will hold the transaction name and we can get it back in TransferTx function in store.go file by calling
//					ctx.Value(txKey) to get the value of txKey from the context. now we have the transaction name and we can add some logs with it
//				./store.go: line 72: txName := ctx.Value(txKey)
//				./store.go: line 73: fmt.Printline(txName, "create transfer")	--- and do the same for the rest of the operations
//				./store.go: line ??: fmt.Printline(txName, "create entry 1")
//				./store.go: line ??: fmt.Printline(txName, "create entry 2")
//				./store.go: line ??: fmt.Printline(txName, "get account 1 for update")
//				./store.go: line ??: fmt.Printline(txName, "update account 1's balance")
//				./store.go: line ??: fmt.Printline(txName, "get account 2 for update")
//				./store.go: line ??: fmt.Printline(txName, "update account 2's balance")
//					but also to make it easier to debug, we should not run too many concurrent transactions and need to
//					change decrease Goroutines from n = 10 to 2. then if we run the test, still we will get deadlock, but this time
//					we have a lot of debuggers there. now we can see what is the cause of this issue, but we need to find why does it happen.
//					for this purpose we can delete all existing data and recreate two accounts and write our transaction query
//					in TablePlus and run it line by line based on the logs we implemented and printed for us:
//				BEGIN;
//				INSERT INTO tranfers (from_account_id, to_account_id, amount) VALUES (1,2,10) RETURNING *;
//				INSERT INTO entries (account_id, amount) VALUES (1,-10) RETURNING *;
//				INSERT INTO entries (account_id, amount) VALUES (2,10) RETURNING *;
//				SELECT * FROM accounts WHERE id = 1 FOR UPDATE;
//				UPDATE accounts SET balance = 90 WHERE id = 1 RETURNING *;
//				SELECT * FROM accounts WHERE id = 2 FOR UPDATE;
//				UPDATE accounts SET balance = 110 WHERE id = 2 RETURNING *;
//				ROLLBACK;
//					then we need to run terminal and with docker command exec running two psql console in paralel to exec both transactions
//					in paralel by running this command: docker exec -it postgres16 psql -U root -d simple_bank
//					then we can see where deacklock appears. to confirm deadlock we can search for 'postgres lock' in Google and find some
//					queries that we can bring them in TablePlus, run and check our process. here after we run and debug queries we can see there
//					is some locks between Accounts table and Transfers table, then if we check our database schema we will see there is only one
//					connection between these two tables which is the account_id as foreign key in Transfers table, and any update on the account_id
//					will affect this foreign key constraint!
//					as we discussed before to solve this issue we can comment those ALTER TABLE statements from our migration file and migratedown
//					and up to remove this connection. but it is not the best way!
//					as we know transaction lock is required because postgres worries that transaction 1 will update account_id, which would affect
//					the foreign key constraints of Transfers table. however if we check the update account query we will see it will only change the
//					balance and update it and account_id will never be changed. so we need to tell postgres that I am selecting this account for update
//					but it's primary key won't be touched then postgres not need a transaction lock, therefore there is no deadlock.
//					to fix it, in GetAccountForUpdate query, instead of SELECT FOR UPDATE we need to tell more clearly SELECT FOR NO KEY UPDATE. this
//					will tell postgress we don't update the key or id column of this table. now we need to regenerate our codes. and then if we run our
//					test again, it will be passed and we should remove debug codes where we add txName and txKey and pass it into the context, and remove
//					remove all logs (fmt.Println()) in transaction codes.
//					still there is a better way to update balance, currently we have to preform 2 queries to get the account and update it's balance.
//					so we can improve this by using a single query to add some amount of money to an account balance directly.
//					for that we're goint to add a new SQL query called AddAcountBalance which is similar to UpdateAccount query. except that, here we
//					set balance equals to balance plus the second argument. if we generate the sqlc query we will se new query is there. however this
//					balance parameter looks a bit confusinig. because we are just adding some amount of money to the balance and not changing the
//					account's balance to this value. so this parameter's name should be amount instead. to do that, in the query instead of using $2
//					for second argument, we can say sqlc.arg(amount) and in the next line instead of $1 we should say WHERE id = sqlc.arg(id).
//					then we need to remove extra codes and instead replace them with new queries we generated.
//			*/
//	}

func TestTransferTx(t *testing.T) {
	// as NewStore function get a new sql.DB object, we need to change our main_test database and add a new global varianle
	//		for it, then pass it into new store
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	// the best practice to handle a database transaction is to run it in several concurrent Go Routines
	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	// it is a wrong way and if we want to check our result and error inside on the below function and
	//		we can not just use testify require functions to check results and errors right there in go routines,
	//		because that function is running inside a different Go Routnine from the one that our TestTransferTx
	// 		function is running on, so there is no guarantee that it will stop the whole test if a condition is not
	//  	satisfied
	// so the correct way to check the result and err is to send them back to the main Goroutine that our test is running on
	//		and send them back there, and that is why we define channels to get connect concurrent Goroutines

	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// we made n transfers in upper lines, now we want to check result
	// check result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// explained why we need this variable in comments below on line 111
		existed := make(map[int]bool)

		// check transfer
		transfer := result.Transfers
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		/*
			till now, we first wrote codes and then unit tests, since now we will do it in TDD
			 so first we wrote our tests with logics in store_test.go, then will come back here
			 to implement logics that pass those tests. so first we check accounts and diff1 and diff2
			 to check they are fine, then after the for loop we should check updated balance for both accounts
		*/

		// check  accounts
		fromAccount := result.FromAccount
		require.NotZero(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount
		// k must be unique for each transaction, 1 for the first transaction, 2 for the second, ...
		//		and n for the nth transaction. in order to check this, we need to declare a new variable
		//		called existed of type map[int]bool
		k := int(diff1 / amount)
		require.True(t, k >= 0 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, updatedAccount1.Balance, account1.Balance-int64(n)*amount)
	require.Equal(t, updatedAccount2.Balance, account2.Balance+int64(n)*amount)
	fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)
}

/*
	in the above unit test we fix the deadlock for the foreign account_id key constraints. but still we have potential deadlock in our

		implementation the best way to do that is to avoid deadlock. so imagine a scenario that account1 send 10USD to account2 and as
		a concurrent transaction, account2 will send 10USD to account1. so new deadlock will happen. to avoid this issue again we will
		follow the TDD and fix this issue with the help of our unit tests and new transfer will implement bellow this old function.
			-- IMPORTANT NOTE -- (IF WE HAVE 2 CONCURRENR TRANSACTIONS INVOLVING THE SAME PAIR OF ACCOUNTS, THERE MIGHT BE A POTENIAL DEADLOCK)
		so to test this scenario we can open two paralel psql console in terminal an follow these transactions queries:

	BEGIN;
	UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
	UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *;
	ROLLBACK;
	BEGIN;
	UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *;
	UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;
	ROLLBACK;

		so if we run the transactions in concurrent manner and first run UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
		then run UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *; and then run UPDATE accounts SET balance = balance
		+ 10 WHERE id = 2 RETURNING *; we will see there is a blocked transaction and if we run the UPDATE accounts SET balance = balance + 10
		WHERE id = 1 RETURNING *; we will get the deadlock. so if we run the queries for postgres we searched from google in previous test we
		will see that there is a deadlock.
		why deadlock happened? because two concurrent transactions try to update the same accounts's balance in different order. with this
		behavior if we run these queries on the below order:

	BEGIN;
	UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
	UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *;
	COMMIT;
	BEGIN;
	UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;
	UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *;
	COMMIT;

		after we run the first UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *; it will complete successfully
		then when we will run UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *; from the second transaction it will lock, now
		when we run UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *; and then COMMIT; the previous lock is unblocked instantly
		and we can run UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *; and then COMMIT;
		to prevent this issue we can easily change our code and let the smaller account_id updates the account first. and that is the reason why we
		added that if/else statement checker for arg.FromAccountID and arg.ToAccountID

.
*/
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	// we want to simulate a scenario that we have 10 transactions, 5 from account1 to account2 and 5 in reverse direction.
	// 		to do so don't need to take care about result and we just want to check error. because we checked it in previous test
	n := 10
	amount := int64(10)
	errs := make(chan error)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID

		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	//	because we have 10 transactions between account1 and account2, which are the same, therefor the final balances are equal
	//		to start balance
	require.Equal(t, updatedAccount1.Balance, account1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)
	fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)
}
