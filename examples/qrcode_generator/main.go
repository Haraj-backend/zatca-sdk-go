package main

import (
	"log"
	"time"

	"github.com/Haraj-backend/zatca-sdk-go/qrcode"
	goqrcode "github.com/skip2/go-qrcode"
)

func main() {
	hash, err := qrcode.EncodeTLV(qrcode.Data{
		SellerName:      "Bobs Records",
		SellerTaxNumber: "310122393500003",
		Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
		InvoiceTotal:    1000,
		TotalVAT:        150,
	})
	if err != nil {
		log.Fatalf("unable to encode TLV due: %v", err)
	}
	err = goqrcode.WriteFile(hash, goqrcode.Medium, 256, "qr.png")
	if err != nil {
		log.Fatalf("unable to encode TLV due: %v", err)
	}
}
