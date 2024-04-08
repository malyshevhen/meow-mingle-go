package db

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testQueries *Queries

func TestCreateUser(t *testing.T) {
	runPostgresContainer(t)

	args := CreateUserParams{
		Email:     "jaja@mail.com",
		FirstName: "Ja Ja",
		LastName:  "Binks",
		Password:  "jajastrong",
	}

	t.Run("test create user with valid params", func(t *testing.T) {
		user, err := testQueries.CreateUser(context.Background(), args)
		require.NoError(t, err)
		require.NotEmpty(t, user)
		require.Equal(t, args.Email, user.Email)
		require.Equal(t, args.FirstName, user.FirstName)
		require.Equal(t, args.LastName, user.LastName)

		require.NotZero(t, user.ID)
		require.NotZero(t, user.CreatedAt)
	})
}

func TestGetUser(t *testing.T) {
	runPostgresContainer(t)

	args := CreateUserParams{
		Email:     "jaja@mail.com",
		FirstName: "Ja Ja",
		LastName:  "Binks",
		Password:  "jajastrong",
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)

	t.Run("test get existing user", func(t *testing.T) {
		user, err = testQueries.GetUser(context.Background(), user.ID)

		require.NoError(t, err)
		require.NotEmpty(t, user)
		require.Equal(t, args.Email, user.Email)
		require.Equal(t, args.FirstName, user.FirstName)
		require.Equal(t, args.LastName, user.LastName)
		require.NotZero(t, user.ID)
		require.NotZero(t, user.CreatedAt)
	})
}

func TestUpdateUser(t *testing.T) {
	runPostgresContainer(t)

	args1 := CreateUserParams{
		Email:     "jaja@mail.com",
		FirstName: "Ja Ja",
		LastName:  "Binks",
		Password:  "jajastrong",
	}

	user1, err := testQueries.CreateUser(context.Background(), args1)
	require.NoError(t, err)

	t.Run("test update existing user", func(t *testing.T) {
		args2 := UpdateUserParams{
			ID:        user1.ID,
			FirstName: "Bob",
			LastName:  "Ross",
		}

		user1, err := testQueries.UpdateUser(context.Background(), args2)
		require.NoError(t, err)

		require.Equal(t, args2.FirstName, user1.FirstName)
		require.Equal(t, args2.LastName, user1.LastName)
	})
}

func TestDeleteUser(t *testing.T) {
	runPostgresContainer(t)

	args := CreateUserParams{
		Email:     "jaja@mail.com",
		FirstName: "Ja Ja",
		LastName:  "Binks",
		Password:  "jajastrong",
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)

	t.Run("test delete existing user", func(t *testing.T) {
		err = testQueries.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)
		user, err = testQueries.GetUser(context.Background(), user.ID)
		require.Error(t, err)
		require.Empty(t, user)
		require.Equal(t, sql.ErrNoRows, err)
		require.Contains(t, err.Error(), "no rows in result set")
	})
}

func TestListUsers(t *testing.T) {
	runPostgresContainer(t)

	args1 := CreateUserParams{
		Email:     "jaja@mail.com",
		FirstName: "Ja Ja",
		LastName:  "Binks",
		Password:  "jajastrong",
	}

	user1, err := testQueries.CreateUser(context.Background(), args1)
	require.NoError(t, err)

	args2 := CreateUserParams{
		Email:     "bob@mail.com",
		FirstName: "Bob",
		LastName:  "Ross",
		Password:  "bobtallent",
	}

	user2, err := testQueries.CreateUser(context.Background(), args2)
	require.NoError(t, err)

	t.Run("test list users", func(t *testing.T) {
		args3 := ListUsersParams{
			Limit:  10,
			Offset: 0,
		}

		users, err := testQueries.ListUsers(context.Background(), args3)

		require.NoError(t, err)
		require.NotEmpty(t, users)
		require.Len(t, users, 2)
		require.Equal(t, user1.ID, users[0].ID)
		require.Equal(t, user2.ID, users[1].ID)
	})
}

func runPostgresContainer(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "migration", "000001_init_schema.up.sql")),
		postgres.WithDatabase("mingle-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("example"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(6*time.Second)),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("can not connect to the DB:", err)
	}

	testQueries = New(conn)
	require.NoError(t, err)
}
