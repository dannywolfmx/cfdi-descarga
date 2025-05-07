package digest

import (
	"fmt"
	"time"

	"github.com/dannywolfmx/cfdi-descarga/signature"
)

type Digest struct {
	Value, Signature string
}

type DigestGenerator struct {
	*signature.Key
	DigestFormat    string
	SignatureFormat string
	CreatedAt       string
	ExpiresAt       string
}

func (d DigestGenerator) New() (*Digest, error) {
	message := fmt.Sprintf(d.DigestFormat, d.CreatedAt, d.ExpiresAt)
	digestValue, err := d.Hash([]byte(message))

	if err != nil {
		return nil, err
	}

	message = fmt.Sprintf(d.SignatureFormat, digestValue)
	signature, err := d.SignMessage(message)

	if err != nil {
		return nil, err
	}

	return &Digest{
		Value:     digestValue,
		Signature: signature,
	}, nil
}

type UUIDV4Generator func() (string, error)

func generateDate(diference int) (string, string) {
	timeFormat := "2006-01-02T15:04:05.000Z"

	t := time.Now()
	creaded := t.UTC().Format(timeFormat)
	expires := t.UTC().Add(time.Minute * time.Duration(diference)).Format(timeFormat)
	return creaded, expires
}
