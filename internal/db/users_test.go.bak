package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	args := CreateUserParams{
		Email:     "jaja@mail.com",
		FirstName: "Ja Ja",
		LastName:  "Binks",
		Password:  "password",
	}

	t.Run("test create user with valid params", func(t *testing.T) {
		teardown := SetupTest(t)
		defer teardown(t)

		user, err := TestUserRepository.CreateUser(context.Background(), args)
		require.NoError(t, err)
		require.NotEmpty(t, user)
		require.Equal(t, args.Email, user.Email)
		require.Equal(t, args.FirstName, user.FirstName)
		require.Equal(t, args.LastName, user.LastName)

		require.NotZero(t, user.ID)
		require.NotZero(t, user.CreatedAt)
	})

	t.Run("test create user with invalid email", func(t *testing.T) {
		teardown := SetupTest(t)
		defer teardown(t)

		invalidEmailArgs := CreateUserParams{
			Email:     "",
			FirstName: "Invalid",
			LastName:  "Email",
			Password:  "password",
		}

		savedUser, err := TestUserRepository.CreateUser(context.Background(), invalidEmailArgs)
		require.Error(t, err)

		t.Logf("Saved user: %v", savedUser)
	})

	t.Run("test create user with duplicate email", func(t *testing.T) {
		teardown := SetupTest(t)
		defer teardown(t)

		_, err := TestUserRepository.CreateUser(context.Background(), args)
		require.NoError(t, err)

		duplicateArgs := CreateUserParams{
			Email:     args.Email,
			FirstName: "Duplicate",
			LastName:  "User",
			Password:  "password",
		}

		_, err = TestUserRepository.CreateUser(context.Background(), duplicateArgs)
		require.Error(t, err)
	})

	t.Run("test create user with invalid password", func(t *testing.T) {
		teardown := SetupTest(t)
		defer teardown(t)

		invalidPasswordArgs := CreateUserParams{
			Email:     "invalid@email.com",
			FirstName: "Invalid",
			LastName:  "Password",
			Password:  "",
		}

		_, err := TestUserRepository.CreateUser(context.Background(), invalidPasswordArgs)
		require.Error(t, err)
	})

	t.Run("test create user with invalid first name", func(t *testing.T) {
		teardown := SetupTest(t)
		defer teardown(t)

		invalidFirstNameArgs := CreateUserParams{
			Email:     "invalidfirstname@email.com",
			FirstName: "",
			LastName:  "Invalid",
			Password:  "First_Name",
		}

		_, err := TestUserRepository.CreateUser(context.Background(), invalidFirstNameArgs)
		require.Error(t, err)
	})

	t.Run("test create user with invalid last name", func(t *testing.T) {
		teardown := SetupTest(t)
		defer teardown(t)

		invalidLastNameArgs := CreateUserParams{
			Email:     "invalidlastname@email.com",
			FirstName: "Invalid",
			LastName:  "",
			Password:  "Last_Name",
		}

		_, err := TestUserRepository.CreateUser(context.Background(), invalidLastNameArgs)
		require.Error(t, err)
	})
}

func TestGetUser(t *testing.T) {
	args := CreateUserParams{
		Email:     "julia@mail.com",
		FirstName: "Julia",
		LastName:  "Roberts",
		Password:  "password",
	}

	t.Run("test get existing user", func(t *testing.T) {
		teardown := SetupTest(t)
		defer teardown(t)

		existingUser, err := TestUserRepository.CreateUser(context.Background(), args)
		require.NoError(t, err)

		user, err := TestUserRepository.GetUser(context.Background(), existingUser.ID)

		require.NoError(t, err)
		require.NotEmpty(t, user)
		require.Equal(t, args.Email, user.Email)
		require.Equal(t, args.FirstName, user.FirstName)
		require.Equal(t, args.LastName, user.LastName)
		require.NotZero(t, user.ID)
		require.NotZero(t, user.CreatedAt)
	})

	t.Run("test error if user does not exists", func(t *testing.T) {
		teardown := SetupTest(t)
		defer teardown(t)

		_, err := TestUserRepository.GetUser(context.Background(), 1_000_000)

		require.Error(t, err)
	})
}
