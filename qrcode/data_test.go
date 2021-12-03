package qrcode_test

import (
	"testing"
	"time"

	"github.com/Haraj-backend/zatca-sdk-go/qrcode"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		Name     string
		Data     qrcode.Data
		ExpError error
	}{
		{
			Name: "Test Valid Data - No Arabic",
			Data: qrcode.Data{
				SellerName:      "Bobs Records",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        150,
			},
			ExpError: nil,
		},
		{
			Name: "Test Valid Data - With Arabic",
			Data: qrcode.Data{
				SellerName:      "الجواهري العربي",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        150,
			},
			ExpError: nil,
		},
		{
			Name: "Test Invalid Data - Long Seller Name Alphabet",
			Data: qrcode.Data{
				SellerName:      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut ultricies sem quis enim pellentesque auctor. In sit amet posuere erat, non aliquam ligula. Vestibulum pretium quis metus vel blandit. Maecenas nec molestie tellus, eget efficitur nunc. Nulla ullamcorper quis nibh eu pretium. Donec et nulla urna. Fusce sapien dolor, consectetur non lorem et, ultrices rhoncus ante. In sed imperdiet mi, sed faucibus massa. Nulla facilisi. Donec molestie eros et metus eleifend faucibus. Nulla scelerisque ex sed turpis.",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        150,
			},
			ExpError: qrcode.ErrSellerNameTooLong,
		},
		{
			Name: "Test Invalid Data - Long Seller Name Arabic",
			Data: qrcode.Data{
				SellerName:      "هنالك العديد من الأنواع المتوفرة لنصوص لوريم إيبسوم، ولكن الغالبية تم تعديلها بشكل ما عبر إدخال بعض النوادر أو الكلمات العشوائية إلى النص. إن كنت تريد أن تستخدم نص لوريم إيبسوم ما، عليك أن تتحقق أولاً أن ليس هناك أي كلمات أو عبارات محرجة أو غير لائقة مخبأة في هذا النص. بينما تعمل جميع مولّدات نصوص لوريم إيبسوم على الإنترنت على إعادة تكرار مقاطع من نص لوريم إيبسوم نفسه عدة مرات بما تتطلبه الحاجة، يقوم مولّدنا هذا باستخدام كلمات من قاموس يحوي على أكثر من 200 كلمة لا تينية، مضاف إليها مجموعة من الجمل النموذجية، لتكوين نص لوريم إيبسوم ذو شكل منطقي قريب إلى النص الحقيقي. وبالتالي يكون النص الناتح خالي من التكرار، أو أي كلمات أو عبارات غير لائقة أو ما شابه. وهذا ما يجعله أول مولّد نص لوريم إيبسوم حقيقي على الإنترنت.",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        150,
			},
			ExpError: qrcode.ErrSellerNameTooLong,
		},
		{
			Name: "Test Invalid Data - Missing Seller Name",
			Data: qrcode.Data{
				SellerName:      "",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        150,
			},
			ExpError: qrcode.ErrMissingSellerName,
		},
		{
			Name: "Test Invalid Data - Missing Seller Tax Number",
			Data: qrcode.Data{
				SellerName:      "Bobs Records",
				SellerTaxNumber: "",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        150,
			},
			ExpError: qrcode.ErrMissingSellerTaxNumber,
		},
		{
			Name: "Test Invalid Data - Missing Invoice Total",
			Data: qrcode.Data{
				SellerName:      "Bobs Records",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    0,
				TotalVAT:        150,
			},
			ExpError: qrcode.ErrMissingInvoiceTotal,
		},
		{
			Name: "Test Invalid Data - Missing Total VAT",
			Data: qrcode.Data{
				SellerName:      "Bobs Records",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        0,
			},
			ExpError: qrcode.ErrMissingTotalVAT,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := testCase.Data.Validate()
			if err != testCase.ExpError {
				t.Errorf("mismatch error, exp: %v, got: %v", testCase.ExpError, err)
			}
		})
	}
}

func TestEncodeTLV(t *testing.T) {
	testCases := []struct {
		Name    string
		Data    qrcode.Data
		ExpHash string
	}{
		{
			Name: "Test No Arabic Name",
			Data: qrcode.Data{
				SellerName:      "Bobs Records",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        150,
			},
			ExpHash: "AQxCb2JzIFJlY29yZHMCDzMxMDEyMjM5MzUwMDAwMwMUMjAyMi0wNC0yNVQxNTozMDowMFoEBzEwMDAuMDAFBjE1MC4wMA==",
		},
		{
			Name: "Test Arabic Name",
			Data: qrcode.Data{
				SellerName:      "الجواهري العربي",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        150,
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
				SellerName:      "Bobs Records",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        150,
			},
		},
		{
			Name:      "Test Arabic Name",
			HashInput: "AR3Yp9mE2KzZiNin2YfYsdmKINin2YTYudix2KjZigIPMzEwMTIyMzkzNTAwMDAzAxQyMDIyLTA0LTI1VDE1OjMwOjAwWgQHMTAwMC4wMAUGMTUwLjAw",
			ExpData: qrcode.Data{
				SellerName:      "الجواهري العربي",
				SellerTaxNumber: "310122393500003",
				Timestamp:       time.Date(2022, 04, 25, 15, 30, 00, 00, time.UTC),
				InvoiceTotal:    1000,
				TotalVAT:        150,
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
