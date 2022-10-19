package piaotong

import (
	"encoding/base64"

	"github.com/sdcxtech/openssl/v2"
)

func (c *Client) encrypt(data []byte) (string, error) {
	encrypted, err := openssl.Des3ECBEncrypt(data, c.tripleDESKey, openssl.PKCS5_PADDING)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (c *Client) Decrypt(content string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}

	decrypted, err := openssl.Des3ECBDecrypt(data, c.tripleDESKey, openssl.PKCS5_PADDING)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
