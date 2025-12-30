package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/idoyudha/duitku-go/common"
)

type APIClient struct {
	*common.ServiceClient
	// API Services
	PaymentService *PaymentService
}

type PaymentService struct {
	client *common.ServiceClient
}

func NewClient(cfg *common.Config) *APIClient {
	c := &APIClient{
		ServiceClient: &common.ServiceClient{
			Cfg: cfg,
		},
	}

	c.PaymentService = &PaymentService{
		client: c.ServiceClient,
	}

	return c
}

var DuitkuClient *APIClient
var serverKey string

func init() {
	if DuitkuClient == nil {
		log.Printf("Initializing Duitku...")
		InitializeMidtrans()
	}
}

func InitializeMidtrans() {
	if DuitkuClient != nil {
		log.Println("Midtrans already initialized")
		return
	}

	DuitkuClient = NewClient(&common.Config{
		Environment:  common.SandboxEnv,
		APIKey:       os.Getenv("DUITKU_API_KEY"),
		MerchantCode: os.Getenv("DUITKU_MERCHANT_CODE"),
	})
}

func (s *PaymentService) setBodyRequest(req *DuitkuRequestCharge) error {
	if s.client.Cfg.MerchantCode == "" {
		return errors.New("merchant code is empty")
	}

	if os.Getenv("DUITKU_CALLBACK_URL") == "" {
		return errors.New("DUITKU_CALLBACK_URL is empty")
	}

	if os.Getenv("DUITKU_RETURN_URL") == "" {
		return errors.New("DUITKU_RETURN_URL is empty")
	}
	req.MerchantCode = s.client.Cfg.MerchantCode
	req.Signature = s.generatePaymentSignature(req.MerchantCode + req.MerchantOrderId + strconv.Itoa(req.PaymentAmount) + s.client.Cfg.APIKey)
	req.CallbackUrl = os.Getenv("DUITKU_CALLBACK_URL")
	req.ReturnUrl = os.Getenv("DUITKU_RETURN_URL")
	return nil
}

func (s *PaymentService) Charge(ctx context.Context, req DuitkuRequestCharge) (DuitkuResponseCharge, *http.Response, error) {
	res := &DuitkuResponseCharge{}
	path := "/merchant/v2/inquiry"

	baseUrl := common.SandboxV2BaseURL
	if s.client.Cfg.Environment == common.ProductionEnv {
		baseUrl = common.ProductionV2BaseURL
	}

	err := s.setBodyRequest(&req)
	if err != nil {
		return *res, nil, err
	}

	httpRes, err := common.SendAPIRequest(
		ctx,
		s.client,
		req,
		res,
		http.MethodPost,
		baseUrl+path,
		nil,
	)

	return *res, httpRes, err
}

func (s *PaymentService) generatePaymentSignature(parameter string) string {
	md5Hash := md5.Sum([]byte(parameter))
	return hex.EncodeToString(md5Hash[:])
}
