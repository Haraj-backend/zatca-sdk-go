package qrcode_test

import (
	"testing"
	"time"

	"github.com/Haraj-backend/zatca-sdk-go/qrcode"
)

func TestValidate(t *testing.T) {}

func TestEncodeTLV(t *testing.T) {
	testCases := []struct {
		Name    string
		Data    qrcode.Data
		ExpHash string
	}{
		{
			Name: "Test No Arabic Name",
			Data: qrcode.Data{
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
			Data: qrcode.Data{
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
			hash, err := qrcode.EncodeTLV(testCase.Data)
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
		Name      string
		HashInput string
		ExpData   qrcode.Data
	}{
		{
			Name:      "Test No Arabic Name",
			HashInput: "AQxCb2JzIFJlY29yZHMCDzMxMDEyMjM5MzUwMDAwMwMUMjAyMi0wNC0yNVQxNTozMDowMFoEBzEwMDAuMDAFBjE1MC4wMA==",
			ExpData: qrcode.Data{
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
			ExpData: qrcode.Data{
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
			if *qrCodeDataResult != testCase.ExpData {
				t.Error("mismatch qrcode data")
			}
		})
	}
}
