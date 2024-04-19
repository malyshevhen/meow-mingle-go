package api

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/malyshEvhen/meow_mingle/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateJwt(t *testing.T) {
	secret := []byte("secret")
	userId := int64(1)

	t.Run("success", func(t *testing.T) {
		token, err := createJwt(secret, userId)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims := &jwt.RegisteredClaims{}
		_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		assert.NoError(t, err)
		assert.Equal(t, userId, claims.Subject)
	})

	t.Run("error signing", func(t *testing.T) {
		_, err := createJwt([]byte("invalid"), userId)

		assert.Error(t, err)
		assert.IsType(t, errors.NewUnauthorizedError(), err)
	})
}

func TestExpiredJwt(t *testing.T) {
	secret := []byte("secret")
	userId := int64(1)

	token, err := createJwt(secret, userId)
	assert.NoError(t, err)

	claims := &jwt.StandardClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	assert.NoError(t, err)

	// Manually override expiry
	claims.ExpiresAt = time.Now().Add(-1 * time.Hour).Unix()

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	assert.Error(t, err)
	assert.IsType(t, errors.NewUnauthorizedError(), err)
}
