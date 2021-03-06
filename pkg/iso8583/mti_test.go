package iso8583_test

import (
	"errors"
	"fmt"
	"github.com/jattento/go-iso8583/pkg/encoding/ebcdic"
	"github.com/jattento/go-iso8583/pkg/iso8583"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMTI_MarshalISO8583(t *testing.T) {
	testList := []struct {
		Name        string
		V           iso8583.MTI
		Encoding    string
		Length      int
		Numba       int
		OutputBytes []byte
		OutputError string
	}{
		{
			Name:        "ascii_standard",
			V:           iso8583.MTI{MTI: "0100"},
			Encoding:    "ascii",
			OutputError: "",
			OutputBytes: []byte("0100"),
		},
		{
			Name:        "ebcdic_standard",
			V:           iso8583.MTI{MTI: "0100"},
			Encoding:    "ebcdic",
			OutputError: "",
			Numba:       0100,
			OutputBytes: ebcdic.V1047.FromGoString("0100"),
		},
	}

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("var_to_bytes_%s", testCase.Name), func(t *testing.T) {
			o, err := testCase.V.MarshalISO8583(testCase.Length, testCase.Encoding)
			if testCase.OutputError != "" {
				assert.Errorf(t, err, testCase.OutputError)
			} else {
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			}
			assert.Equal(t, testCase.OutputBytes, o)
		})
	}
}

func TestMTI_UnmarshalISO8583(t *testing.T) {
	testList := []struct {
		Name          string
		InputBytes    []byte
		InputEncoding string
		InputLength   int
		OutputContent string
		OutputError   string
	}{
		{
			Name:          "ascii_standard",
			InputEncoding: "ascii",
			InputLength:   4,
			OutputContent: "0110",
			OutputError:   "",
			InputBytes:    []byte("0110"),
		},
		{
			Name:          "ebcdic_standard",
			InputEncoding: "ebcdic",
			InputLength:   4,
			OutputContent: "0110",
			OutputError:   "",
			InputBytes:    ebcdic.V1047.FromGoString("0110"),
		},
		{
			Name:          "error_length",
			InputEncoding: "ebcdic",
			InputLength:   5,
			OutputContent: "",
			OutputError:   "mti isnt 4 characters long, its: 5",
			InputBytes:    ebcdic.V1047.FromGoString("01100"),
		},
		{
			Name:          "error_text",
			InputEncoding: "ebcdic",
			InputLength:   4,
			OutputContent: "",
			OutputError:   "mti characters arent numbers: strconv.Atoi: parsing \"text\": invalid syntax",
			InputBytes:    ebcdic.V1047.FromGoString("text"),
		},
		{
			Name:          "error_text",
			InputEncoding: "ebcdic",
			InputLength:   4,
			OutputContent: "",
			OutputError:   "mti characters arent numbers: strconv.Atoi: parsing \"text\": invalid syntax",
			InputBytes:    ebcdic.V1047.FromGoString("text"),
		},
		{
			Name:          "unmarshal_error",
			InputEncoding: "force_error",
			InputLength:   4,
			OutputContent: "",
			OutputError:   "encoder 'force_error' returned error: forced_error",
			InputBytes:    []byte{1, 1, 1, 1, 1},
		},
		{
			Name:          "len_shorter_than_length_error",
			InputEncoding: "ascii",
			InputLength:   4,
			OutputContent: "",
			OutputError:   "message remain (1 bytes) is shorter than indicated length: 4",
			InputBytes:    []byte("1"),
		},
		{
			Name:          "nil_bytes_error",
			InputEncoding: "ascii",
			InputLength:   4,
			OutputContent: "",
			OutputError:   "bytes input is nil",
			InputBytes:    nil,
		},
	}

	iso8583.UnmarshalDecodings["force_error"] = func(bytes []byte) ([]byte, error) { return nil, errors.New("forced_error") }
	defer delete(iso8583.UnmarshalDecodings, "force_error")

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("var_to_bytes_%s", testCase.Name), func(t *testing.T) {
			var v iso8583.MTI

			_, err := v.UnmarshalISO8583(testCase.InputBytes, testCase.InputLength, testCase.InputEncoding)
			if testCase.OutputError != "" {
				if assert.NotNil(t, err) {
					assert.Equal(t, err.Error(), testCase.OutputError)
				}
				return
			} else {
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			}
			assert.Equal(t, testCase.OutputContent, v.String())
		})
	}
}
