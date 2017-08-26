package serkis

import "encoding/base64"

var cryptoKey = NewEncryptionKey()

type Share struct {
	Fpath string
}

func NewShareFromSecret(secret string, key *[32]byte) (Share, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return Share{}, err
	}

	text, err := Decrypt(ciphertext, key)
	if err != nil {
		return Share{}, err
	}

	return Share{Fpath: string(text)}, nil
}

func (s Share) Secret(key *[32]byte) (string, error) {
	ciphertext, err := Encrypt([]byte(s.Fpath), key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
