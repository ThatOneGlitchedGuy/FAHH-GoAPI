package security

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/fernet/fernet-go"
)

func Encrypt(text []byte, keyString string) ([]byte, error) {
	key, err := base64.URLEncoding.DecodeString(keyString)
	if err != nil {
		return nil, errors.New("invalid fernet key format")
	}
	if len(key) != 32 {
		return nil, errors.New("fernet key must be 32 bytes")
	}

	k, err := fernet.DecodeKey(keyString)
	if err != nil {
		return nil, err
	}

	encrypted, err := fernet.EncryptAndSign(text, k)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

func Decrypt(encrypted []byte, keyString string, ttl time.Duration) ([]byte, error) {
	key, err := base64.URLEncoding.DecodeString(keyString)
	if err != nil {
		return nil, errors.New("invalid fernet key format")
	}
	if len(key) != 32 {
		return nil, errors.New("fernet key must be 32 bytes")
	}

	k, err := fernet.DecodeKey(keyString)
	if err != nil {
		return nil, err
	}

	decrypted := fernet.VerifyAndDecrypt(encrypted, ttl, []*fernet.Key{k})
	if decrypted == nil {
		return nil, errors.New("failed to verify or decrypt message")
	}
	return decrypted, nil
}
