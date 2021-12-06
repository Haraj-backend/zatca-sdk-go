# Zatca SDK GO

[![Build](https://github.com/Haraj-backend/zatca-sdk-go/actions/workflows/build.yml/badge.svg)](https://github.com/Haraj-backend/zatca-sdk-go/actions/workflows/build.yml)
[![Test](https://github.com/Haraj-backend/zatca-sdk-go/actions/workflows/test.yml/badge.svg)](https://github.com/Haraj-backend/zatca-sdk-go/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Haraj-backend/zatca-sdk-go)](https://goreportcard.com/report/github.com/Haraj-backend/zatca-sdk-go)
[![GitHub License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/user/project/master/LICENSE)

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
	fmt.Println("hash data:", hash)

	// decode hash using TLV method to get data
	data, err := qrcode.DecodeTLV(hash)
	if err != nil {
		log.Fatalf("unable to decode TLV due: %v", err)
	}
	fmt.Printf("decoded hash: %s", data)
}
```

## Generating QR Code

This package is only used for encoding QR Code data into base64 hash using TLV method. So it doesn't contain functionality to generate QR Code image.

If you want to generate the QR Code image, you could use another library such as [skip2/go-qrcode](https://github.com/skip2/go-qrcode). You could check [examples/qrcode_generator](./examples/qrcode_generator) for details.

## Other Projects
Here the list about similar projects that using other programming languages:
1. [Node (axenda/zatca)](https://github.com/axenda/zatca)
2. [Ruby (mrsool/zatca)](https://github.com/mrsool/zatca)
3. [Python (TheAwiteb/fatoora)](https://github.com/TheAwiteb/fatoora)
4. [Kotlin & Java (iabdelgawaad/ZATCA)](https://github.com/iabdelgawaad/ZATCA)
5. [Swift (elgawady14/ZATCA)](https://github.com/elgawady14/ZATCA)
6. Javascript
	- [husninazer/fatoora-ksa](https://github.com/husninazer/fatoora-ksa)
	- [Evincible-Solutions/EVSZatcaQRCodeJavascript](https://github.com/Evincible-Solutions/EVSZatcaQRCodeJavascript)
7. DotNet / DotNetCore
	- [aljbri/Zatca.Net](https://github.com/aljbri/Zatca.Net)
	- [alquhait/ZatcaDotNetCore](https://github.com/alquhait/ZatcaDotNetCore)
	- [Evincible-Solutions/EvsZatcaQRCodeString](https://github.com/Evincible-Solutions/EvsZatcaQRCodeString)
8. PHP
	- [SallaApp/ZATCA](https://github.com/SallaApp/ZATCA)
	- [MukhtarSayedSaleh/saudi-zakat-qr-generator](https://github.com/MukhtarSayedSaleh/saudi-zakat-qr-generator)
	- [mPhpMaster/laravel-zatca](https://github.com/mPhpMaster/laravel-zatca)
	- [IdaraNet/ZATCA-PHP-TLV-QR-CODE](https://github.com/IdaraNet/ZATCA-PHP-TLV-QR-CODE)
9. [REST API (NafieAlhilaly/api-fatoora)](https://github.com/NafieAlhilaly/api-fatoora)

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

1. If you have any suggestions to make this SDK better, feel free to fork this project and create a pull request
2. If you found any bugs, you can report [here](https://github.com/Haraj-backend/zatca-sdk-go/issues)

## License
The MIT License (MIT). Please see [License File](LICENSE) for more information.
