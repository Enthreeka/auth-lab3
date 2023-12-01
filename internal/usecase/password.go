package usecase

import (
	"crypto/subtle"
	"encoding/hex"
	"github.com/Enthreeka/auth-lab3/internal/apperror"
	"golang.org/x/crypto/argon2"
)

type argon struct {
	salt    []byte
	version int
	threads uint8
	time    uint32
	memory  uint32
	keyLen  uint32
}

// TODO Create global salt
// Salt settings for hash generation. Salt = user_id + ...
func (a *argon) setSalt(salt string) {
	userSalt := []byte(salt)
	a.salt = userSalt
}

func NewPassword(salt string) *argon {
	return &argon{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}
}

func (a *argon) GenerateHashFromPassword(id string, password string) (string, error) {
	//if !validation.IsValidEmail(user.Login) && !validation.IsValidPassword(user.Password) {
	//	return "", apperror.ErrDataNotValid
	//}

	a.setSalt(id)

	hashPasswordByte := argon2.IDKey([]byte(password), a.salt, a.time, a.memory, a.threads, a.keyLen)
	hashPasswordString := hex.EncodeToString(hashPasswordByte)

	return hashPasswordString, nil
}

func (a *argon) VerifyPassword(hashPassword string, id string, password string) error {
	a.setSalt(id)

	newHashByte := argon2.IDKey([]byte(password), a.salt, a.time, a.memory, a.threads, a.keyLen)

	newHashString := hex.EncodeToString(newHashByte)

	if subtle.ConstantTimeCompare([]byte(hashPassword), []byte(newHashString)) != 1 {
		return apperror.ErrHashPasswordsNotEqual
	}

	return nil
}
