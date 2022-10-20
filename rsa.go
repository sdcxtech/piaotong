package piaotong

import (
	"crypto"
	"fmt"

	"github.com/sdcxtech/openssl/v2"
)

// SignRequest 对向票通发起的请求进行签名
func (c *Client) SignRequest(req *Request) error {
	signed, err := openssl.RSASign([]byte(req.SignatureContent()), c.platformPrivateKey, crypto.SHA1)
	if err != nil {
		return err
	}

	req.Sign = base64EncodeToString(signed)

	return nil
}

// SignResponse 对返回给票通的响应进行签名
func (c *Client) SignResponse(res *Response) error {
	signed, err := openssl.RSASign([]byte(res.SignatureContent()), c.platformPrivateKey, crypto.SHA1)
	if err != nil {
		return err
	}

	res.Sign = base64EncodeToString(signed)

	return nil
}

// VerifyRequest 验证票通请求签名
func (c *Client) VerifyRequest(req *Request) error {
	sign, err := base64DecodeString(req.Sign)
	if err != nil {
		return err
	}

	err = openssl.RSAVerify([]byte(req.SignatureContent()), sign, c.piaotongPublicKey, crypto.SHA1)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrInvalidSignature, err)
	}

	return err
}

// VerifyResponse 验证票通响应签名
func (c *Client) VerifyResponse(res *Response) error {
	sign, err := base64DecodeString(res.Sign)
	if err != nil {
		return err
	}

	err = openssl.RSAVerify([]byte(res.SignatureContent()), sign, c.piaotongPublicKey, crypto.SHA1)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrInvalidSignature, err)
	}

	return err
}
