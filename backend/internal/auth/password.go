package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

var ErrInvalidPasswordHash = errors.New("invalid password hash")

type PasswordHashParams struct {
	MemoryKB    uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type PasswordHasher struct {
	params PasswordHashParams
}

func NewPasswordHasher(params PasswordHashParams) PasswordHasher {
	return PasswordHasher{params: params}
}

func DefaultPasswordHasher() PasswordHasher {
	return NewPasswordHasher(PasswordHashParams{
		MemoryKB:    64 * 1024,
		Iterations:  1,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	})
}

func (hasher PasswordHasher) Hash(password string) (string, error) {
	salt := make([]byte, hasher.params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(password), salt, hasher.params.Iterations, hasher.params.MemoryKB, hasher.params.Parallelism, hasher.params.KeyLength)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		hasher.params.MemoryKB,
		hasher.params.Iterations,
		hasher.params.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func (hasher PasswordHasher) Verify(password string, encodedHash string) bool {
	params, salt, expectedKey, err := parsePasswordHash(encodedHash)
	if err != nil {
		return false
	}

	actualKey := argon2.IDKey([]byte(password), salt, params.Iterations, params.MemoryKB, params.Parallelism, uint32(len(expectedKey)))

	return subtle.ConstantTimeCompare(actualKey, expectedKey) == 1
}

func parsePasswordHash(encodedHash string) (PasswordHashParams, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return PasswordHashParams{}, nil, nil, ErrInvalidPasswordHash
	}

	params := PasswordHashParams{}
	for _, pair := range strings.Split(parts[3], ",") {
		keyValue := strings.SplitN(pair, "=", 2)
		if len(keyValue) != 2 {
			return PasswordHashParams{}, nil, nil, ErrInvalidPasswordHash
		}

		value, err := strconv.ParseUint(keyValue[1], 10, 32)
		if err != nil {
			return PasswordHashParams{}, nil, nil, ErrInvalidPasswordHash
		}

		switch keyValue[0] {
		case "m":
			params.MemoryKB = uint32(value)
		case "t":
			params.Iterations = uint32(value)
		case "p":
			params.Parallelism = uint8(value)
		default:
			return PasswordHashParams{}, nil, nil, ErrInvalidPasswordHash
		}
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return PasswordHashParams{}, nil, nil, ErrInvalidPasswordHash
	}

	key, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return PasswordHashParams{}, nil, nil, ErrInvalidPasswordHash
	}

	params.SaltLength = uint32(len(salt))
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}
