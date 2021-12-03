# Zatca SDK GO

[![Build](https://github.com/Haraj-backend/zatca-sdk-go/actions/workflows/build.yml/badge.svg)](https://github.com/Haraj-backend/zatca-sdk-go/actions/workflows/build.yml)
[![Test](https://github.com/Haraj-backend/zatca-sdk-go/actions/workflows/test.yml/badge.svg)](https://github.com/Haraj-backend/zatca-sdk-go/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Haraj-backend/zatca-sdk-go)](https://goreportcard.com/report/github.com/Haraj-backend/zatca-sdk-go)

An unofficial package in Golang to help developers to implement ZATCA (Fatoora) QR code easily which required for e-invoicing

> ✅ The hash result has been validated the same as the output from ZATCA's SDK as of 18th November 2021

## Installation

```
go get github.com/Haraj-backend/zatca-sdk-go
```

## Simple Usage

```golang
package main

import (
	"log"
	"time"

	"github.com/Haraj-backend/zatca-sdk-go/qrcode"
)

func main() {
	// encode data using TLV method to get code hash
	hash, err := qrCode.EncodeTLV(qrcode.Data{
		SellerName:   		"Bobs Records",
		SellerTaxNumber:    "310122393500003",
		Timestamp:    		time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
		InvoiceTotal: 		1000,
		TotalVAT:     		150,
	})
	if err != nil {
		log.Fatalf("unable to encode TLV due: %v", err)
	}
	fmt.Println("hash data:", hash)

	// decode hash using TLV method to get data
	data, err := qrcode.DecodeTLV("AR3Yp9mE2KzZiNin2YfYsdmKINin2YTYudix2KjZigIPMzEwMTIyMzkzNTAwMDAzAxQyMDIyLTA0LTI1VDE1OjMwOjAwWgQHMTAwMC4wMAUGMTUwLjAw")
	if err != nil {
		log.Fatalf("unable to decode TLV due: %v", err)
	}
	fmt.Printf("decoded hash: %s", data)
}
```

## Generating QR Code

This package is only used for encoding QR Code data into base64 hash using TLV method. So it doesn't contain functionality to generate QR Code image.

If you want to generate the QR Code image, you could use another library such as [skip2/go-qrcode](https://github.com/skip2/go-qrcode). You could check [examples/qrcode_generator](./examples/qrcode_generator) for details.
