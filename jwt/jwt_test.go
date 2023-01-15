package jwt

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestJwtSecretKey(t *testing.T) {
	t.Run("Should return JWT_KEY env variable", func(t *testing.T) {
		envKey := "ultra_secret_key"
		os.Setenv("JWT_KEY", envKey)

		secret := jwtSecretKey()
		expected := []byte(envKey)
		if !bytes.Equal(secret, expected) {
			t.Errorf("JWT_KEY (%s), expected: (%s)", string(secret), string(expected))
		}
	})
}

func TestJwtExpirationTime(t *testing.T) {
	t.Run("Should use default expiration time when missing JWT_EXPIRATION_TIME env var", func(t *testing.T) {
		expirationT := jwtExpirationTime()

		if expirationT.Hour() != time.Now().Hour()+1 {
			t.Errorf("Default expiration time should be 1 hour")
		}
	})

	t.Run("Should use JWT_EXPIRATION_TIME env var", func(t *testing.T) {
		envKey := "120"
		os.Setenv("JWT_EXPIRATION_TIME", envKey)

		expirationT := jwtExpirationTime()

		if expirationT.Hour() != time.Now().Hour()+2 {
			t.Errorf("Expiration time should be 2 hour")
		}
	})
}

func TestNewToken(t *testing.T) {
	t.Run("should token as string", func(t *testing.T) {
		envKey := "ultra_secret_key"
		os.Setenv("JWT_KEY", envKey)

		_, err := NewToken()
		if err != nil {
			t.Errorf("Could not create token, cause: %s", err.Error())
		}
	})
}
