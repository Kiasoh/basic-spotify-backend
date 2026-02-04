package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Password encapsulates the logic for password hashing and comparison.
type Password struct {
	Plaintext *string
	Hash      []byte
}

// Set computes the bcrypt hash of a plaintext password.
func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash

	return nil
}

// Matches checks if a plaintext password matches a given hash.
func (p *Password) Matches(plaintextPassword string, hashedPassword []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(plaintextPassword))
	if err != nil {
		switch {
		case err == bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// Bytes returns the password hash as a byte slice.
func (p *Password) Bytes() []byte {
    return p.Hash
}

// String returns the password hash as a string. Note that this is not reversible.
func (p *Password) String() string {
    return string(p.Hash)
}

func IsValidPassword(password string) (bool, error) {
	if len(password) < 8 {
		return false, fmt.Errorf("password must be at least 8 characters long")
	}

	return true, nil
}
