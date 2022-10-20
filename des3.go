package piaotong

import (
	"github.com/forgoer/openssl"
)

func (c *Client) Encrypt(data []byte) (string, error) {
	encrypted, err := openssl.Des3ECBEncrypt(data, c.tripleDESKey, openssl.PKCS5_PADDING)
	if err != nil {
		return "", err
	}

	return base64EncodeToString(encrypted), nil
}

func (c *Client) Decrypt(content string) ([]byte, error) {
	data, err := base64DecodeString(content)
	if err != nil {
		return nil, err
	}

	decrypted, err := openssl.Des3ECBDecrypt(data, c.tripleDESKey, openssl.PKCS5_PADDING)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
