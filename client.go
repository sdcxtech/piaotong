package piaotong

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jkomyno/nanoid"
)

var ErrInvalidSignature = errors.New("invalid signature")

var timezoneBeijing = time.FixedZone("Beijing Time", int(8*time.Hour.Seconds()))

type Client struct {
	platformCode       string
	platformName       string
	platformPrivateKey []byte
	piaotongPublicKey  []byte
	tripleDESKey       []byte

	restyClient *resty.Client
}

type Config struct {
	Host               string
	PlatformCode       string
	PlatformName       string // 平台简称
	PlatformPrivateKey string // PKCS #1, PEM form
	PiaotongPublicKey  string // PKCS #1, PEM form
	TripleDESKey       string
}

func New(config Config) *Client {
	restyClient := resty.New().
		SetBaseURL(config.Host)

	return &Client{
		platformCode:       config.PlatformCode,
		platformName:       config.PlatformName,
		platformPrivateKey: []byte(config.PlatformPrivateKey),
		piaotongPublicKey:  []byte(config.PiaotongPublicKey),
		tripleDESKey:       []byte(config.TripleDESKey),
		restyClient:        restyClient,
	}
}

func (c *Client) PlatformCode() string {
	return c.platformCode
}

func (c *Client) PlatformName() string {
	return c.platformName
}

func (c *Client) GenerateSerialNo() (string, error) {
	id, err := nanoid.Nanoid(8)
	if err != nil {
		return "", err
	}

	ts := time.Now().In(timezoneBeijing).Format("20060102150405")

	return strings.Join([]string{c.platformName, ts, id}, ""), nil
}

func (c *Client) GenerateInvoiceReqSerialNo() (string, error) {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	id, err := nanoid.Generate(alphabet, 16)
	if err != nil {
		return "", err
	}

	return c.platformName + id, nil
}

type (
	OpenBlueInvoiceRequest struct {
		TaxpayerNum        string `json:"taxpayerNum"`
		InvoiceReqSerialNo string `json:"invoiceReqSerialNo"`
		BuyerName          string `json:"buyerName"`
		BuyerTaxpayerNum   string `json:"buyerTaxpayerNum,omitempty"`
		BuyerAddress       string `json:"buyerAddress,omitempty"`
		BuyerTel           string `json:"buyerTel,omitempty"`
		BuyerBankName      string `json:"buyerBankName,omitempty"`
		BuyerBankAccount   string `json:"buyerBankAccount,omitempty"`
		SellerAddress      string `json:"sellerAddress,omitempty"`
		SellerTel          string `json:"sellerTel,omitempty"`
		SellerBankName     string `json:"sellerBankName,omitempty"`
		SellerBankAccount  string `json:"sellerBankAccount,omitempty"`
		ItemName           string `json:"itemName,omitempty"`
		CasherName         string `json:"casherName,omitempty"`
		ReviewerName       string `json:"reviewerName,omitempty"`
		DrawerName         string `json:"drawerName,omitempty"`
		TakerName          string `json:"takerName,omitempty"`
		TakerTel           string `json:"takerTel,omitempty"`
		TakerEmail         string `json:"takerEmail,omitempty"`
		SpecialInvoiceKind string `json:"specialInvoiceKind,omitempty"`
		Remark             string `json:"remark,omitempty"`
		DefinedData        string `json:"definedData,omitempty"`
		TradeNo            string `json:"tradeNo,omitempty"`
		ExtensionNum       string `json:"extensionNum,omitempty"`
		MachineCode        string `json:"machineCode,omitempty"`
		AgentInvoiceFlag   string `json:"agentInvoiceFlag,omitempty"`
		ShopNum            string `json:"shopNum,omitempty"`

		ItemList []*GoodsItem `json:"itemList"`
	}

	GoodsItem struct {
		GoodsName              string `json:"goodsName"`
		TaxClassificationCode  string `json:"taxClassificationCode"`
		SpecificationModel     string `json:"specificationModel,omitempty"`
		MeteringUnit           string `json:"meteringUnit,omitempty"`
		Quantity               string `json:"quantity,omitempty"`
		IncludeTaxFlag         string `json:"includeTaxFlag,omitempty"`
		UnitPrice              string `json:"unitPrice,omitempty"`
		InvoiceAmount          string `json:"invoiceAmount"`
		TaxRateValue           string `json:"taxRateValue"`
		TaxRateAmount          string `json:"taxRateAmount,omitempty"`
		DiscountAmount         string `json:"discountAmount,omitempty"`
		DiscountTaxRateAmount  string `json:"discountTaxRateAmount,omitempty"`
		DeductionAmount        string `json:"deductionAmount,omitempty"`
		PreferentialPolicyFlag string `json:"preferentialPolicyFlag,omitempty"`
		ZeroTaxFlag            string `json:"zeroTaxFlag,omitempty"`
		VatSpecialManage       string `json:"vatSpecialManage,omitempty"`
	}

	OpenBlueInvoiceResponse struct {
		InvoiceReqSerialNo string `json:"invoiceReqSerialNo"`
		QRCodePath         string `json:"qrCodePath"` // base64 decoded
		QRCode             string `json:"qrCode"`     // base64 encoded
	}
)

func (c *Client) OpenBlueInvoice(
	ctx context.Context, req *OpenBlueInvoiceRequest,
) (*OpenBlueInvoiceResponse, error) {
	resp := &OpenBlueInvoiceResponse{}

	err := c.request(ctx, "/tp/openapi/invoiceBlue.pt", req, resp)
	if err != nil {
		return nil, err
	}

	qrcodePath, err := base64DecodeString(resp.QRCodePath)
	if err != nil {
		return nil, err
	}

	resp.QRCodePath = string(qrcodePath)

	return resp, nil
}

type (
	OpenRedInvoiceRequest struct {
		TaxpayerNum        string `json:"taxpayerNum"`
		InvoiceReqSerialNo string `json:"invoiceReqSerialNo"`
		InvoiceCode        string `json:"invoiceCode"`
		InvoiceNo          string `json:"invoiceNo"`
		RedReason          string `json:"redReason"`
		Amount             string `json:"amount"`
		DefinedData        string `json:"definedData,omitempty"`
	}

	OpenRedInvoiceResponse struct {
		InvoiceReqSerialNo string `json:"invoiceReqSerialNo"`
		QRCodePath         string `json:"qrCodePath"`       // base64 decoded
		QRCode             string `json:"qrCode,omitempty"` // base64 encoded
	}
)

func (c *Client) OpenRedInvoice(
	ctx context.Context, req *OpenRedInvoiceRequest,
) (*OpenRedInvoiceResponse, error) {
	resp := &OpenRedInvoiceResponse{}

	err := c.request(ctx, "/tp/openapi/invoiceRed.pt", req, resp)
	if err != nil {
		return nil, err
	}

	qrcodePath, err := base64DecodeString(resp.QRCodePath)
	if err != nil {
		return nil, err
	}

	resp.QRCodePath = string(qrcodePath)

	return resp, nil
}

type (
	QueryInvoiceRequest struct {
		TaxpayerNum        string `json:"taxpayerNum"`
		InvoiceReqSerialNo string `json:"invoiceReqSerialNo"`
	}

	QueryInvoiceResponse struct {
		TaxpayerNum              string `json:"taxpayerNum"`
		InvoiceReqSerialNo       string `json:"invoiceReqSerialNo"`
		InvoiceType              string `json:"invoiceType"`
		Code                     string `json:"code"`
		Msg                      string `json:"msg"`
		TradeNo                  string `json:"tradeNo,omitempty"`
		SecurityCode             string `json:"securityCode,omitempty"`
		QRCode                   string `json:"qrCode,omitempty"`
		InvoiceCode              string `json:"invoiceCode,omitempty"`
		InvoiceNo                string `json:"invoiceNo,omitempty"`
		InvoiceDate              string `json:"invoiceDate,omitempty"`
		NoTaxAmount              string `json:"noTaxAmount,omitempty"`
		TaxAmount                string `json:"taxAmount,omitempty"`
		InvoiceLayoutFileType    string `json:"invoiceLayoutFileType,omitempty"`
		InvoicePdf               string `json:"invoicePdf,omitempty"`
		DownloadURL              string `json:"downloadUrl,omitempty"`              // base64 decoded
		VatPlatformInvPreviewURL string `json:"vatPlatformInvPreviewUrl,omitempty"` // base64 decoded
		ExtensionNum             string `json:"extensionNum,omitempty"`
		DiskNo                   string `json:"diskNo,omitempty"`
		InvPreviewQrcodePath     string `json:"invPreviewQrcodePath,omitempty"` // base64 decoded
		InvPreviewQrcode         string `json:"invPreviewQrcode,omitempty"`     // base64 encoded
	}
)

func (c *Client) QueryInvoice(
	ctx context.Context, req *QueryInvoiceRequest,
) (*QueryInvoiceResponse, error) {
	resp := &QueryInvoiceResponse{}

	err := c.request(ctx, "/tp/openapi/queryInvoice.pt", req, resp)
	if err != nil {
		return nil, err
	}

	downloadURL, err := base64DecodeString(resp.DownloadURL)
	if err != nil {
		return nil, err
	}
	resp.DownloadURL = string(downloadURL)

	vatPlatformInvPreviewURL, err := base64DecodeString(resp.VatPlatformInvPreviewURL)
	if err != nil {
		return nil, err
	}
	resp.VatPlatformInvPreviewURL = string(vatPlatformInvPreviewURL)

	invPreviewQrcodePath, err := base64DecodeString(resp.InvPreviewQrcodePath)
	if err != nil {
		return nil, err
	}
	resp.InvPreviewQrcodePath = string(invPreviewQrcodePath)

	return resp, nil
}

type (
	RegisterRequest struct {
		TaxpayerNum                string `json:"taxpayerNum"`
		EnterpriseName             string `json:"enterpriseName"`
		LegalPersonName            string `json:"legalPersonName"`
		ContactsName               string `json:"contactsName"`
		ContactsEmail              string `json:"contactsEmail"`
		ContactsPhone              string `json:"contactsPhone"`
		RegionCode                 string `json:"regionCode"`
		CityName                   string `json:"cityName"`
		EnterpriseAddress          string `json:"enterpriseAddress"`
		TaxRegistrationCertificate string `json:"taxRegistrationCertificate"`
		TaxControlDeviceType       string `json:"taxControlDeviceType"`
	}

	RegisterResponse struct {
		TaxpayerNum    string `json:"taxpayerNum"`
		EnterpriseName string `json:"enterpriseName"`
	}
)

func (c *Client) Register(
	ctx context.Context, req *RegisterRequest,
) (*RegisterResponse, error) {
	resp := &RegisterResponse{}

	err := c.request(ctx, "/tp/openapi/register.pt", req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type (
	GetEnterpriseInfoRequest struct {
		TaxpayerNum    string `json:"taxpayerNum"`
		EnterpriseName string `json:"enterpriseName"`
	}

	GetEnterpriseInfoResponse struct {
		TaxpayerNum               string `json:"taxpayerNum"`
		EnterpriseName            string `json:"enterpriseName"`
		LegalPersonName           string `json:"legalPersonName"`
		ContactsName              string `json:"contactsName"`
		ContactsEmail             string `json:"contactsEmail"`
		ContactsPhone             string `json:"contactsPhone"`
		RegionCode                string `json:"regionCode"`
		CityName                  string `json:"cityName"`
		EnterpriseAddress         string `json:"enterpriseAddress"`
		InvitationCode            string `json:"invitationCode"`
		ReviewStatus              string `json:"reviewStatus"`
		ReviewOpinion             string `json:"reviewOpinion"`
		TerminalType              string `json:"terminalType"`
		InvoiceKind               string `json:"invoiceKind"`
		InvoiceLayoutFileType     string `json:"invoiceLayoutFileType"`
		BlockchainInvSingleQuota  string `json:"blockchainInvSingleQuota"`
		BlockchainInvDailyQuota   string `json:"blockchainInvDailyQuota"`
		BlockchainInvMonthlyQuota string `json:"blockchainInvMonthlyQuota"`
		ServiceStatus             string `json:"serviceStatus"`

		TerminalList []*TerminalItem `json:"terminalList"`
	}

	TerminalItem struct {
		DiskType         string `json:"diskType"`
		ExtensionNum     string `json:"extensionNum"`
		MachineCode      string `json:"machineCode"`
		Available        string `json:"available"`
		ServiceStartTime string `json:"serviceStartTime"`
		ServiceEndTime   string `json:"serviceEndTime"`
	}
)

func (c *Client) GetEnterpriseInfo(
	ctx context.Context, req *GetEnterpriseInfoRequest,
) (*GetEnterpriseInfoResponse, error) {
	resp := &GetEnterpriseInfoResponse{}

	err := c.request(ctx, "/tp/openapi/getEnterpriseInfo.pt", req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) request(ctx context.Context, url string, reqContent, respPtr any) error {
	reqBody, err := c.buildRequest(reqContent)
	if err != nil {
		return err
	}

	respBody := &Response{}
	err = c.post(ctx, url, reqBody, respBody)
	if err != nil {
		return err
	}

	err = c.VerifyResponse(respBody)
	if err != nil {
		return err
	}

	if respBody.isError() {
		return respBody.toError()
	}

	err = c.decryptAndUnmarshalResponse(respBody.Content, respPtr)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) buildRequest(content any) (*Request, error) {
	data, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	encrypted, err := c.Encrypt(data)
	if err != nil {
		return nil, err
	}

	serialNo, err := c.GenerateSerialNo()
	if err != nil {
		return nil, err
	}

	req := &Request{
		PlatformCode: c.platformCode,
		SignType:     "RSA",
		Format:       "JSON",
		Timestamp:    time.Now().In(timezoneBeijing).Format("2006-01-02 15:04:05"),
		Version:      "1.0",
		SerialNo:     serialNo,
		Content:      encrypted,
	}
	err = c.SignRequest(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) decryptAndUnmarshalResponse(content string, v any) error {
	data, err := c.Decrypt(content)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) post(ctx context.Context, url string, body, resultPtr any) error {
	req := c.restyClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		ForceContentType("application/json")

	resp, err := req.SetBody(body).Post(url)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("http status: %s, body: %s", resp.Status(), resp.Body())
	}

	if resultPtr != nil {
		err = json.Unmarshal(resp.Body(), resultPtr)
		if err != nil {
			return err
		}
	}

	return nil
}

func base64EncodeToString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func base64DecodeString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
