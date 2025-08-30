// Package security contains everything related to the security of user data
package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type ArgonHash struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func New() *ArgonHash {
	return &ArgonHash{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func (a *ArgonHash) GenerateFromPassword(p string) (encoded string, err error) {
	salt, err := genRandByt(a.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(p), salt, a.Iterations, a.Memory, a.Parallelism, a.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded = fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		a.Memory, a.Iterations, a.Parallelism, b64Salt, b64Hash)

	return encoded, nil
}

// VerifyPasswd compares a password p with the stored PHC-style encoded hash e
func (a *ArgonHash) VerifyPasswd(p, e string) (ok bool, err error) {
	parts := strings.Split(e, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var memory, iterations uint32
	var parallelism uint8

	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	calcHash := argon2.IDKey([]byte(p), salt, iterations, memory, parallelism, uint32(len(hash)))

	return subtle.ConstantTimeCompare(hash, calcHash) == 1, nil
}

func genRandByt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
