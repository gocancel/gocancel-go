package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultTolerance indicates that signatures older than this will be rejected.
	DefaultTolerance time.Duration = 300 * time.Second
	// signingVersion represents the version of the signature we currently use.
	signingVersion string = "v1"
)

// This block represents the list of errors that could be raised when using the webhook package.
var (
	ErrInvalidHeader    = errors.New("webhook has invalid Gocxl-Signature header")
	ErrNoValidSignature = errors.New("webhook had no valid signature")
	ErrNotSigned        = errors.New("webhook has no Gocxl-Signature header")
	ErrTooOld           = errors.New("timestamp wasn't within tolerance")
)

// ComputeSignature computes a webhook signature using GoCancel's v1 signing method.
func ComputeSignature(t time.Time, payload []byte, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d", t.Unix())))
	mac.Write([]byte("."))
	mac.Write(payload)

	return mac.Sum(nil)
}

// ValidatePayload validates the payload against the Gocxl-Signature header
// using the specified signing secret. Returns an error if the body or
// Gocxl-Signature header provided are unreadable, if the signature doesn't
// match, or if the timestamp for the signature is older than DefaultTolerance.
func ValidatePayload(payload []byte, header string, secret string) error {
	return ValidatePayloadWithTolerance(payload, header, secret, DefaultTolerance)
}

// ValidatePayloadIgnoringTolerance validates the payload against the Gocxl-Signature header
// header using the specified signing secret. Returns an error if the body or
// Gocxl-Signature header provided are unreadable or if the signature doesn't match.
// Does not check the signature's timestamp.
func ValidatePayloadIgnoringTolerance(payload []byte, header string, secret string) error {
	return validatePayload(payload, header, secret, 0*time.Second, false)
}

// ValidatePayloadWithTolerance validates the payload against the Gocxl-Signature header
// using the specified signing secret and tolerance window. Returns an error if the body
// or Gocxl-Signature header provided are unreadable, if the signature doesn't match, or
// if the timestamp for the signature is older than the specified tolerance.
func ValidatePayloadWithTolerance(payload []byte, header string, secret string, tolerance time.Duration) error {
	return validatePayload(payload, header, secret, tolerance, true)
}

type signedHeader struct {
	timestamp  time.Time
	signatures [][]byte
}

func parseSignatureHeader(header string) (*signedHeader, error) {
	sh := &signedHeader{}

	if header == "" {
		return sh, ErrNotSigned
	}

	// Signed header looks like "t=1495999758,v1=ABC,v1=DEF,v0=GHI"
	pairs := strings.Split(header, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			return sh, ErrInvalidHeader
		}

		switch parts[0] {
		case "t":
			timestamp, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return sh, ErrInvalidHeader
			}
			sh.timestamp = time.Unix(timestamp, 0)

		case signingVersion:
			sig, err := hex.DecodeString(parts[1])
			if err != nil {
				continue // Ignore invalid signatures
			}

			sh.signatures = append(sh.signatures, sig)

		default:
			continue // Ignore unknown parts of the header
		}
	}

	if len(sh.signatures) == 0 {
		return sh, ErrNoValidSignature
	}

	return sh, nil
}

func validatePayload(payload []byte, sigHeader string, secret string, tolerance time.Duration, enforceTolerance bool) error {
	header, err := parseSignatureHeader(sigHeader)
	if err != nil {
		return err
	}

	expectedSignature := ComputeSignature(header.timestamp, payload, secret)
	expiredTimestamp := time.Since(header.timestamp) > tolerance
	if enforceTolerance && expiredTimestamp {
		return ErrTooOld
	}

	// Check all given v1 signatures, multiple signatures will be sent temporarily in the case of a rolled signature secret
	for _, sig := range header.signatures {
		if hmac.Equal(expectedSignature, sig) {
			return nil
		}
	}

	return ErrNoValidSignature
}
