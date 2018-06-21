package anondb

import (
	"encoding/base64"
	"errors"

	"github.com/satori/go.uuid"
)

func generateUUID() string {
	return base64urlEncode(uuid.NewV4().Bytes())
}

// base64urlEncode implements an URL safe base64 encoding, as seen here:
// https://tools.ietf.org/html/rfc4648#section-5
func base64urlEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// ErrNotFound signals that the queried object was not found in the database
var ErrNotFound = errors.New("Item not found")

// ErrDuplicate signals that the insertion would have been a duplicate
var ErrDuplicate = errors.New("Item is duplicate")
