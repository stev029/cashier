package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type PaymentMethod string

const (
	SPAYQR    PaymentMethod = "SP"
	NOBUQR    PaymentMethod = "NQ"
	GUDANGVQR PaymentMethod = "GQ"
	NUSAQR    PaymentMethod = "SQ"
	QRIS      PaymentMethod = SPAYQR
)

type DuitkuResponseCharge struct {
	MerchantCode  string `json:"merchantCode"`
	Reference     string `json:"reference"`
	Amount        string `json:"amount"`
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
	QrString      string `json:"qrString"`
}

type DuitkuRequestCharge struct {
	MerchantCode    string        `json:"merchantCode"`
	PaymentAmount   int           `json:"paymentAmount"`
	MerchantOrderId string        `json:"merchantOrderId"`
	ProductDetails  string        `json:"productDetails"`
	Email           string        `json:"email"`
	PaymentMethod   PaymentMethod `json:"paymentMethod"`
	CustomerVaName  string        `json:"customerVaName"`
	ReturnUrl       string        `json:"returnUrl"`
	CallbackUrl     string        `json:"callbackUrl"`
	Signature       string        `json:"signature"`
}

func SendAPIRequest(
	ctx context.Context,
	req any,
	res any,
	method string,
	url string,
	headerParams map[string]string,
) (*http.Response, error) {
	r, err := setRequest(ctx, method, url, req, headerParams)
	if err != nil {
		return nil, err
	}

	httpResp, err := http.DefaultClient.Do(r)
	if err != nil {
		return httpResp, err
	}

	body, err := io.ReadAll(httpResp.Body)
	httpResp.Body.Close()
	httpResp.Body = io.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		return httpResp, err
	}

	if err = json.Unmarshal(body, &res); err != nil {
		return httpResp, err
	}

	return httpResp, nil
}

func setRequest(
	ctx context.Context,
	method string,
	urlInput string,
	reqBody any,
	headerParams map[string]string,
) (req *http.Request, err error) {
	var body *bytes.Buffer

	if reqBody != nil {
		body = &bytes.Buffer{}
		err = json.NewEncoder(body).Encode(reqBody)
		if err != nil {
			return nil, err
		}
	}

	parsedUrl, err := url.Parse(urlInput)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req, err = http.NewRequestWithContext(ctx, method, parsedUrl.String(), body)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, parsedUrl.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers.Set(h, v)
		}
		req.Header = headers
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
