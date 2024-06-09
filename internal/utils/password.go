package utils

import (
	"log/slog"

	"github.com/alexedwards/argon2id"
)

func GenerateHashFromPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		slog.Error("error in hashing password", err)
		return "", err
	}
	return hash, nil
}

func ComparePasswordAndHash(password string, encodedHash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, encodedHash)
	if err != nil {
		slog.Error("error in compairing password", err)
		return false, err
	}
	return match, nil
}
