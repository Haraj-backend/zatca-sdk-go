package main

import (
	"log"
	"time"

	"github.com/Haraj-backend/zatca-sdk-go/qrcode"
	goqrcode "github.com/skip2/go-qrcode"
)

func main() {
	input := qrcode.Data{
		SellerName:   "Bobs Records",
		SellerTRN:    "310122393500003",
		Timestamp:    time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
		InvoiceTotal: 1000,
		TotalVAT:     150,
	}

	qrCode, err := qrcode.NewQRCode(input)
	if err != nil {
		log.Fatalf("unable to initialize new qrcode due: %v", err)
	}

	hash, err := qrCode.EncodeTLV()
	if err != nil {
		log.Fatalf("unable to encode TLV due: %v", err)
	}

	err = goqrcode.WriteFile(hash, goqrcode.Medium, 256, "qr.png")
	if err != nil {
		log.Fatalf("unable to encode TLV due: %v", err)
	}
}
