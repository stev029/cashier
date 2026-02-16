package utils

import (
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

var MidtransClient coreapi.Client

func InitializeMidtrans() {
	MidtransClient = coreapi.Client{}
	midtransEnv := midtrans.Sandbox
	if os.Getenv("MIDTRANS_ENV") == "production" {
		midtransEnv = midtrans.Production
	}

	MidtransClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtransEnv)
}
