package user

import (
	"app/pkg/validator"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	Plaintext    *string `json:"password"`
	Confirmation *string `json:"password_confirmation"`
	Previous     *string `json:"previous_password"`
	Hash         *[]byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = &hash
	return nil
}

func (p *password) Match(hash *[]byte) (bool, error) {
	if p.Plaintext == nil || hash == nil {
		return false, nil
	}
	if err := bcrypt.CompareHashAndPassword(
		*hash,
		[]byte(*p.Plaintext),
	); err != nil {
		return false, err
	}
	return true, nil
}

func (p *password) CheckPreviousPassword(v *validator.Validator, userID string) error {
	var hashedPassword []byte

	query := "SELECT password_hash FROM users WHERE id = $1"
	err := v.DB.GetContext(context.Background(), &hashedPassword, query, userID)
	if err != nil {
		return err
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(*p.Previous))
	if err != nil {
		// Passwords do not match
		return err
	}

	// Passwords match
	return nil
}
