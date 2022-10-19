package piaotong

import (
	"crypto"
	"encoding/base64"

	"github.com/sdcxtech/openssl/v2"
)

func (c *Client) sign(s string) (string, error) {
	signed, err := openssl.RSASign([]byte(s), c.platformPrivateKey, crypto.SHA1)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signed), nil
}

func (c *Client) Verify(req *Request) error {
	sign, err := base64.StdEncoding.DecodeString(req.Sign)
	if err != nil {
		return err
	}

	s := c.buildSignatureContent(req)

	return openssl.RSAVerify([]byte(s), sign, c.piaotongPublicKey, crypto.SHA1)
}
