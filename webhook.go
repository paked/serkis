package serkis

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"strings"
)

var (
	NoSignatureError                   = errors.New("No signature")
	IncorrectlyFormattedSignatureError = errors.New("Invalid formatted signature")
)

type WebhookValidator struct {
	Secret    string
	Signature string
	Body      io.Reader
}

func (wv *WebhookValidator) Valid() (bool, error) {
	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(wv.Signature) == 0 {
		return false, NoSignatureError
	}

	if len(wv.Signature) != signatureLength || !strings.HasPrefix(wv.Signature, signaturePrefix) {
		return false, IncorrectlyFormattedSignatureError
	}

	body, err := ioutil.ReadAll(wv.Body)
	if err != nil {
		return false, err
	}

	computed := hmac.New(sha1.New, []byte(wv.Secret))
	computed.Write(body)
	generated := []byte(computed.Sum(nil))

	return hmac.Equal(generated, wv.Actual()), nil
}

func (wv *WebhookValidator) Actual() []byte {
	actual := make([]byte, 20)
	hex.Decode(actual, []byte(wv.Signature[5:]))

	return actual
}
