package qrcode_test

import (
	"testing"
	"time"

	"github.com/Haraj-backend/zatca-sdk-go/qrcode"
)

func TestEncodeTLV(t *testing.T) {
	testCases := []struct {
		Name       string
		QRCodeData qrcode.QRCodeData
		ExpHash    string
	}{
		{
			Name: "Test No Arabic Name",
			QRCodeData: qrcode.QRCodeData{
				SellerName:   "Bobs Records",
				SellerTRN:    "310122393500003",
				Timestamp:    time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal: 1000,
				TotalVAT:     150,
			},
			ExpHash: "AQxCb2JzIFJlY29yZHMCDzMxMDEyMjM5MzUwMDAwMwMUMjAyMi0wNC0yNVQxNTozMDowMFoEBzEwMDAuMDAFBjE1MC4wMA==",
		},
		{
			Name: "Test Arabic Name",
			QRCodeData: qrcode.QRCodeData{
				SellerName:   "الجواهري العربي",
				SellerTRN:    "310122393500003",
				Timestamp:    time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal: 1000,
				TotalVAT:     150,
			},
			ExpHash: "AR3Yp9mE2KzZiNin2YfYsdmKINin2YTYudix2KjZigIPMzEwMTIyMzkzNTAwMDAzAxQyMDIyLTA0LTI1VDE1OjMwOjAwWgQHMTAwMC4wMAUGMTUwLjAw",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			qrCode, err := qrcode.NewQRCode(testCase.QRCodeData)
			if err != nil {
				t.Fatalf("unable to initialize new qrcode due: %v", err)
			}
			hash, err := qrCode.EncodeTLV()
			if err != nil {
				t.Fatalf("unable to get hash due: %v", err)
			}
			if testCase.ExpHash != hash {
				t.Error("mismatch hash")
			}
		})
	}
}

func TestDecodeTLV(t *testing.T) {
	testCases := []struct {
		Name               string
		HashInput          string
		ExpectedQRCodeData qrcode.QRCodeData
	}{
		{
			Name:      "Test No Arabic Name",
			HashInput: "AQxCb2JzIFJlY29yZHMCDzMxMDEyMjM5MzUwMDAwMwMUMjAyMi0wNC0yNVQxNTozMDowMFoEBzEwMDAuMDAFBjE1MC4wMA==",
			ExpectedQRCodeData: qrcode.QRCodeData{
				SellerName:   "Bobs Records",
				SellerTRN:    "310122393500003",
				Timestamp:    time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal: 1000,
				TotalVAT:     150,
			},
		},
		{
			Name:      "Test Arabic Name",
			HashInput: "AR3Yp9mE2KzZiNin2YfYsdmKINin2YTYudix2KjZigIPMzEwMTIyMzkzNTAwMDAzAxQyMDIyLTA0LTI1VDE1OjMwOjAwWgQHMTAwMC4wMAUGMTUwLjAw",
			ExpectedQRCodeData: qrcode.QRCodeData{
				SellerName:   "الجواهري العربي",
				SellerTRN:    "310122393500003",
				Timestamp:    time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal: 1000,
				TotalVAT:     150,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			qrCodeDataResult, err := qrcode.DecodeTLV(testCase.HashInput)
			if err != nil {
				t.Fatalf("unable to decode TLV to qrcode data due: %v", err)
			}

			if *qrCodeDataResult != testCase.ExpectedQRCodeData {
				t.Error("mismatch qrcode data")
			}
		})
	}
}
