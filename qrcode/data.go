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

// maxValueLength define the maximum length for the field value inside Data
// since the length could only be 1 byte, that means the maximum length for
// every field values is 255.
const maxValueLength = 255

// Data represents decoded data inside QRCode
type Data struct {
	SellerName   string
	SellerTRN    string
	Timestamp    time.Time
	InvoiceTotal float64
	TotalVAT     float64
}

// Validate is used for validating Data
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

// String returns string representation of Data
func (d Data) String() string {
	out, _ := json.Marshal(d)
	return string(out)
}

// EncodeTLV return base64 hash value for given Data. Internally this function will
// call `Data.Validate()` for validating the input data.
func EncodeTLV(d Data) (string, error) {
	// validate data
	err := d.Validate()
	if err != nil {
		return "", err
	}
	// construct hash
	builder := &strings.Builder{}
	builder.WriteString(encodeValue(1, d.SellerName))
	builder.WriteString(encodeValue(2, d.SellerTRN))
	builder.WriteString(encodeValue(3, formatTime(d.Timestamp)))
	builder.WriteString(encodeValue(4, formatFloat(d.InvoiceTotal)))
	builder.WriteString(encodeValue(5, formatFloat(d.TotalVAT)))

	b, err := hex.DecodeString(builder.String())
	if err != nil {
		return "", fmt.Errorf("unable to decode hex string due: %v", err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// DecodeTLV returns Data for given base64 hash string.
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

		err = decodeValue(idx, val, &data)
		if err != nil {
			return nil, fmt.Errorf("unable to set value for idx: %d val: %s due: %v", idx, val, err)
		}
		tempBytes := bytesData[2+length:]
		bytesData = tempBytes
	}

	return &data, nil
}

func decodeValue(idx int, val string, d *Data) error {
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

func encodeValue(idx int, val string) string {
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
