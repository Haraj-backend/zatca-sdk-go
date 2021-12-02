package qrcode

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const maxValueLength = 255

type QRCode interface {
	EncodeTLV() (string, error)
}

type QRCodeData struct {
	SellerName   string
	SellerTRN    string
	Timestamp    time.Time
	InvoiceTotal float64
	TotalVAT     float64
}

func (i QRCodeData) ValidateInput() error {
	if len(i.SellerName) == 0 {
		return fmt.Errorf("missing `SellerName`")
	}
	if len(i.SellerTRN) == 0 {
		return fmt.Errorf("missing `SellerTRN`")
	}
	if i.Timestamp == (time.Time{}) {
		return fmt.Errorf("missing `Timestamp`")
	}
	if i.InvoiceTotal == 0 {
		return fmt.Errorf("missing `InvoiceTotal`")
	}
	if i.TotalVAT == 0 {
		return fmt.Errorf("missing `TotalVAT")
	}
	strMap := map[string]string{
		"SellerName":   i.SellerName,
		"SellerTRN":    i.SellerTRN,
		"Timestamp":    formatTime(i.Timestamp),
		"InvoiceTotal": formatFloat(i.InvoiceTotal),
		"TotalVAT":     formatFloat(i.TotalVAT),
	}
	for fieldName, str := range strMap {
		if len([]byte(str)) > maxValueLength {
			return fmt.Errorf("`%v` exceeding max value length", fieldName)
		}
	}
	return nil
}

// SetValue to set qrcode value based on index number
// it follows Zatca QR Code specification
func (i *QRCodeData) SetValue(idx int, val string) error {
	var err error
	switch idx {
	case 1:
		i.SellerName = val
	case 2:
		i.SellerTRN = val
	case 3:
		i.Timestamp, err = time.Parse(time.RFC3339, val)
		if err != nil {
			return fmt.Errorf("timestamp format shoud be in RFC3339")
		}
	case 4:
		i.InvoiceTotal, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("invoice total format shoud be float number")
		}
	case 5:
		i.TotalVAT, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("total VAT format shoud be float number")
		}
	}
	return nil
}

func NewQRCode(input QRCodeData) (QRCode, error) {
	err := input.ValidateInput()
	if err != nil {
		return nil, err
	}
	c := &qrCode{
		sellerName:   input.SellerName,
		sellerTRN:    input.SellerTRN,
		timestamp:    input.Timestamp,
		invoiceTotal: input.InvoiceTotal,
		totalVAT:     input.TotalVAT,
	}
	return c, nil
}

type qrCode struct {
	sellerName   string
	sellerTRN    string
	timestamp    time.Time
	invoiceTotal float64
	totalVAT     float64
}

func (d qrCode) EncodeTLV() (string, error) {
	builder := &strings.Builder{}
	builder.WriteString(constructTLV(1, d.sellerName))
	builder.WriteString(constructTLV(2, d.sellerTRN))
	builder.WriteString(constructTLV(3, formatTime(d.timestamp)))
	builder.WriteString(constructTLV(4, formatFloat(d.invoiceTotal)))
	builder.WriteString(constructTLV(5, formatFloat(d.totalVAT)))

	b, err := hex.DecodeString(builder.String())
	if err != nil {
		return "", fmt.Errorf("unable to decode hex string due: %v", err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeTLV(hash string) (*QRCodeData, error) {
	bytesData, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64: %v", err)
	}

	qrCodeData := QRCodeData{}
	// read TLV data
	for len(bytesData) > 0 {
		// get index from byte stands for TAG in TLV format
		idx := int(bytesData[0])
		// get length value from byte stands for LENGTH in TLV format
		length := bytesData[1]
		// get value from bytes stands for VALUE in TLV format
		// convert the bytes to string
		val := string(bytesData[2 : 2+length])

		err = qrCodeData.SetValue(idx, val)
		if err != nil {
			return nil, fmt.Errorf("unable to set value for idx: %d val: %s due: %v", idx, val, err)
		}
		tempBytes := bytesData[2+length:]
		bytesData = tempBytes
	}

	return &qrCodeData, nil
}

func constructTLV(idx int, val string) string {
	builder := &strings.Builder{}
	rns := []byte(val)
	for i := 0; i < len(rns); i++ {
		builder.WriteString(fmt.Sprintf("%x", rns[i]))
	}
	return fmt.Sprintf("%02x%02x%s", idx, len(rns), builder)
}

func formatTime(t time.Time) string {
	return t.In(time.UTC).Format(time.RFC3339)
}

func formatFloat(x float64) string {
	return fmt.Sprintf("%.2f", x)
}
