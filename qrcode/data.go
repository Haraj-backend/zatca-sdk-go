package qrcode

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// maxValueLength define the maximum length for the field value inside Data
// since the length could only be 1 byte, that means the maximum length for
// every field values is 255.
const maxValueLength = 255

var (
	ErrMissingSellerName      = errors.New("missing `SellerName`")
	ErrMissingSellerTaxNumber = errors.New("missing `SellerTaxNumber`")
	ErrMissingTimestamp       = errors.New("missing `Timestamp`")
	ErrMissingInvoiceTotal    = errors.New("missing `InvoiceTotal`")
	ErrMissingTotalVAT        = errors.New("missing `TotalVAT`")

	ErrSellerNameTooLong      = errors.New("`SellerName` is too long")
	ErrSellerTaxNumberTooLong = errors.New("`SellerTaxNumber` is too long")
)

// Data represents decoded data inside QRCode
type Data struct {
	SellerName      string
	SellerTaxNumber string
	Timestamp       time.Time
	InvoiceTotal    float64
	TotalVAT        float64
}

// Validate is used for validating Data
func (i Data) Validate() error {
	if len(i.SellerName) == 0 {
		return ErrMissingSellerName
	}
	if len(i.SellerTaxNumber) == 0 {
		return ErrMissingSellerTaxNumber
	}
	if i.Timestamp == (time.Time{}) {
		return ErrMissingTimestamp
	}
	if i.InvoiceTotal == 0 {
		return ErrMissingInvoiceTotal
	}
	if i.TotalVAT == 0 {
		return ErrMissingTotalVAT
	}
	strMap := map[string]string{
		"SellerName":      i.SellerName,
		"SellerTaxNumber": i.SellerTaxNumber,
	}
	errMap := map[string]error{
		"SellerName":      ErrSellerNameTooLong,
		"SellerTaxNumber": ErrSellerTaxNumberTooLong,
	}
	for fieldName, str := range strMap {
		if len(str) > maxValueLength {
			return errMap[fieldName]
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
	buf := new(bytes.Buffer)
	buf.Write(encodeValue(1, d.SellerName))
	buf.Write(encodeValue(2, d.SellerTaxNumber))
	buf.Write(encodeValue(3, formatTime(d.Timestamp)))
	buf.Write(encodeValue(4, formatFloat(d.InvoiceTotal)))
	buf.Write(encodeValue(5, formatFloat(d.TotalVAT)))

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
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

		err = setValue(idx, val, &data)
		if err != nil {
			return nil, fmt.Errorf("unable to set value for idx: %d val: %s due: %v", idx, val, err)
		}
		tempBytes := bytesData[2+length:]
		bytesData = tempBytes
	}

	return &data, nil
}

func setValue(idx int, val string, d *Data) error {
	var err error
	switch idx {
	case 1:
		d.SellerName = val
	case 2:
		d.SellerTaxNumber = val
	case 3:
		d.Timestamp, err = time.Parse(time.RFC3339, val)
		if err != nil {
			return fmt.Errorf("timestamp format should be in RFC3339")
		}
	case 4:
		d.InvoiceTotal, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("invoice total format should be float number")
		}
	case 5:
		d.TotalVAT, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("total VAT format should be float number")
		}
	}
	return nil
}

func encodeValue(idx int, val string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(idx))      // write `Tag`
	buf.WriteByte(byte(len(val))) // write `Length`
	buf.Write([]byte(val))        // write `Value`

	return buf.Bytes()
}

func formatTime(t time.Time) string {
	return t.In(time.UTC).Format(time.RFC3339)
}

func formatFloat(x float64) string {
	return fmt.Sprintf("%.2f", x)
}
