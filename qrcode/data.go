package qrcode

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const maxValueLength = 255

type Data struct {
	SellerName   string
	SellerTRN    string
	Timestamp    time.Time
	InvoiceTotal float64
	TotalVAT     float64
}

func (i Data) Validate() error {
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

func (d Data) String() string {
	out, _ := json.Marshal(d)
	return string(out)
}

func EncodeTLV(d Data) (string, error) {
	// validate data
	err := d.Validate()
	if err != nil {
		return "", err
	}
	// construct hash
	builder := &strings.Builder{}
	builder.WriteString(constructTLV(1, d.SellerName))
	builder.WriteString(constructTLV(2, d.SellerTRN))
	builder.WriteString(constructTLV(3, formatTime(d.Timestamp)))
	builder.WriteString(constructTLV(4, formatFloat(d.InvoiceTotal)))
	builder.WriteString(constructTLV(5, formatFloat(d.TotalVAT)))

	b, err := hex.DecodeString(builder.String())
	if err != nil {
		return "", fmt.Errorf("unable to decode hex string due: %v", err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeTLV(hash string) (*Data, error) {
	bytesData, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64: %v", err)
	}

	data := Data{}
	// read TLV data
	for len(bytesData) > 0 {
		// get index from byte stands for `Tag` in TLV format
		idx := int(bytesData[0])
		// get length value from byte stands for `Length` in TLV format
		length := bytesData[1]
		// get value from bytes stands for `Value` in TLV format
		// convert the bytes to string
		val := string(bytesData[2 : 2+length])

		err = data.decodeValue(idx, val)
		if err != nil {
			return nil, fmt.Errorf("unable to set value for idx: %d val: %s due: %v", idx, val, err)
		}
		tempBytes := bytesData[2+length:]
		bytesData = tempBytes
	}

	return &data, nil
}

func (d *Data) decodeValue(idx int, val string) error {
	var err error
	switch idx {
	case 1:
		d.SellerName = val
	case 2:
		d.SellerTRN = val
	case 3:
		d.Timestamp, err = time.Parse(time.RFC3339, val)
		if err != nil {
			return fmt.Errorf("timestamp format shoud be in RFC3339")
		}
	case 4:
		d.InvoiceTotal, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("invoice total format shoud be float number")
		}
	case 5:
		d.TotalVAT, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("total VAT format shoud be float number")
		}
	}
	return nil
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
