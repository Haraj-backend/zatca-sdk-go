# Zatca SDK GO

An unofficial package in Golang to help developers to implement ZATCA (Fatoora) QR code easily which required for e-invoicing

# Installation

```
go get github.com/Haraj-backend/zatca-sdk-go
```

# Simple Usage
```golang
package main

import (
	"log"
	"time"

	"github.com/Haraj-backend/zatca-sdk-go/qrcode"
)

func main() {
	qrCode, err := qrcode.NewQRCode(qrcode.QRCodeData{
		SellerName:   "Bobs Records",
		SellerTRN:    "310122393500003",
		Timestamp:    time.Now(),
		InvoiceTotal: 1000,
		TotalVAT:     150,
	})
	if err != nil {
		log.Fatalf("unable to initialize new qrcode due: %v", err)
	}

	// encode TLV to get hash data
	hash, err := qrCode.EncodeTLV()
	if err != nil {
		log.Fatalf("unable to encode TLV due: %v", err)
	}
	fmt.Println("hash data:", hash)

	// decode TLV to get QR Code data
	qrCodeResult, err := qrcode.DecodeTLV("AR3Yp9mE2KzZiNin2YfYsdmKINin2YTYudix2KjZigIPMzEwMTIyMzkzNTAwMDAzAxQyMDIyLTA0LTI1VDE1OjMwOjAwWgQHMTAwMC4wMAUGMTUwLjAw")
	if err != nil {
		log.Fatalf("unable to decode TLV due: %v", err)
	}
	fmt.Println("decoded hash:", qrCodeResult)
}
```
