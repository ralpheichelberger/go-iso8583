# go-iso8583

An easy to use, yet flexible marshaler for ISO-8583

This library provides:
- Marshal and Unmarshal functions with his respective interfaces
including MTI, VAR, LLVAR, LLLVAR and bitmaps fields ready for use
but with the possibility to easily add new field types.
- Inbuid Support for ASCII, BCD and EBCDIC but not limited to them

## Installation

To install go-iso8583 package, you need to install Go and set your Go workspace first.

1. First you need [Go](https://golang.org/) installed, then you can use the below Go command to install go-iso8583.
```sh
$ go get -u github.com/jattento/go-iso8583
```

2. Import it in your code:
```go
import "github.com/jattento/go-iso8583/pkg/iso8583"
```

## Quick start

```go
import "github.com/go-iso8583/pkg/iso8583"

type PurchaseRequest struct {
	MTI                    iso8583.MTI    `iso8583:"mti"`
	FirstBitmap            iso8583.BITMAP `iso8583:"bitmap,length:64"` // length is the maximum amount of represented elements.
	SecondBitmap           iso8583.BITMAP `iso8583:"1,length:64"`      // length is the maximum amount of represented elements.
	PAN                    iso8583.LLVAR  `iso8583:"2"`
	ProcessingCode         iso8583.VAR    `iso8583:"3"`
	Amount                 iso8583.VAR    `iso8583:"4,encoding:ebcdic"` // By default ASCII is assumed but dont limit yourself!
	DateTime               iso8583.VAR    `iso8583:"7"`
	SystemTraceAuditNumber iso8583.VAR    `iso8583:"11,omitempty"` // omitempty is supported!
	LocalTransactionTime   iso8583.VAR    `iso8583:"12"`
	LocalTransactionDate   iso8583.VAR    `iso8583:"-"` // You can explicitly ignore a field.
	ExpirationDate         iso8583.VAR    `iso8583:"14"`
	MerchantType           iso8583.VAR    `iso8583:"18"`
	ICC                    iso8583.LLLVAR `iso8583:"55"`
	SettlementCode         iso8583.VAR    `iso8583:"66"`
	MessageNumber          iso8583.VAR    `iso8583:"71"`
	TransactionDescriptor  iso8583.VAR    `iso8583:"104"`
}

func GenerateStaticReqBytes() ([]byte, error) {
	req := PurchaseRequest{
		MTI: "0100",
		// FirstBitmap is generated by library
		// SecondBitmap is generated by library
		PAN: "54000000000000111", // LL part is added by library!
		ProcessingCode: "1000",
		Amount: "000000000100",
		MessageNumber: "1",
	}
	
	byt, err := iso8583.Marshal(req)
	if err != nil {
		return nil, err
	}

	return byt, nil
}
```

```go
import "github.com/go-iso8583/pkg/iso8583"

type PurchaseResponse struct {
	MTI                    iso8583.MTI    `iso8583:"mti,length:4"`
	FirstBitmap            iso8583.BITMAP `iso8583:"bitmap,length:64"` // length is the maximum amount of represented elements.
	SecondBitmap           iso8583.BITMAP `iso8583:"1,length:64"`      // length is the maximum amount of represented elements.
	PAN                    iso8583.LLVAR  `iso8583:"2,length:2"`       // length is the amount of bytes of the LL part.
	ProcessingCode         iso8583.VAR    `iso8583:"3,length:6"`
	Amount                 iso8583.VAR    `iso8583:"4,length:12"`
	DateTime               iso8583.VAR    `iso8583:"7,length:10"`
	SystemTraceAuditNumber iso8583.VAR    `iso8583:"11,length:6"`
	LocalTransactionTime   iso8583.VAR    `iso8583:"12,length:6"`
	LocalTransactionDate   iso8583.VAR    `iso8583:"13,length:4"`
	ExpirationDate         iso8583.VAR    `iso8583:"14,length:4"`
	MerchantType           iso8583.VAR    `iso8583:"18,length:4"`
	ResponseCode           iso8583.VAR    `iso8583:"39,length:2"`
	ICC                    iso8583.LLLVAR `iso8583:"55,length:3,encoding:ebcdic/ascii"` // LLL and VAR part use different encoding? Use a / to indicate both
	SettlementCode         iso8583.VAR    `iso8583:"66,length:1"`
	MessageNumber          iso8583.VAR    `iso8583:"71,length:4"`
	TransactionDescriptor  iso8583.VAR    `iso8583:"104,length:100"`
}

func ReadResp(byt []byte) (PurchaseResponse,error){
	var resp PurchaseResponse

	_, err :=iso8583.Unmarshal(byt,&resp)
	if err != nil{
		return PurchaseResponse{}, err
	}

	return resp,nil
}
```