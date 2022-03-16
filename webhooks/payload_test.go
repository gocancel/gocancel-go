package webhooks

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

var (
	testPayload = []byte(`{
		"id": "evt_test_webhook",
		"object": "event"
  	}`)
	testSecret = "wh_sig_test_secret"
)

type signedPayload struct {
	timestamp time.Time
	payload   []byte
	secret    string
	scheme    string
	signature []byte
	header    string
}

func generateHeader(p signedPayload) string {
	return fmt.Sprintf("t=%d,%s=%s", p.timestamp.Unix(), p.scheme, hex.EncodeToString(p.signature))
}

func newSignedPayload(options ...func(*signedPayload)) *signedPayload {
	signedPayload := &signedPayload{}
	signedPayload.timestamp = time.Now()
	signedPayload.payload = testPayload
	signedPayload.secret = testSecret
	signedPayload.scheme = "v1"

	for _, opt := range options {
		opt(signedPayload)
	}

	if signedPayload.signature == nil {
		signedPayload.signature = ComputeSignature(signedPayload.timestamp, signedPayload.payload, signedPayload.secret)
	}
	signedPayload.header = generateHeader(*signedPayload)
	return signedPayload
}

func (p *signedPayload) hexSignature() string {
	return hex.EncodeToString(p.signature)
}

func TestValidatePayload(t *testing.T) {
	p := newSignedPayload()
	err := ValidatePayload(p.payload, "", p.secret)
	if err != ErrNotSigned {
		t.Errorf("Expected ErrNotSigned from missing signature, got %v", err)
	}

	err = ValidatePayload(p.payload, "t=", p.secret)
	if err != ErrInvalidHeader {
		t.Errorf("Expected ErrInvalidHeader from bad header format, got %v", err)
	}

	err = ValidatePayload(p.payload, p.header+",v1=bad_signature", p.secret)
	if err != nil {
		t.Errorf("Received unexpected %v error with an unreadable signature in the header (should be ignored)", err)
	}

	p = newSignedPayload(func(p *signedPayload) {
		p.scheme = "v0"
	})
	err = ValidatePayload(p.payload, p.header, p.secret)
	if err != ErrNoValidSignature {
		t.Errorf("Expected error from mismatched schema, got %v", err)
	}

	p = newSignedPayload(func(p *signedPayload) {
		p.signature = []byte("deadbeef")
	})
	err = ValidatePayload(p.payload, p.header, p.secret)
	if err != ErrNoValidSignature {
		t.Errorf("Expected error from fake signature, got %v", err)
	}

	p = newSignedPayload()
	p2 := newSignedPayload(func(p *signedPayload) {
		p.secret = testSecret + "_rolled_key"
	})
	headerWithRolledKey := p.header + ",v1=" + p2.hexSignature()
	if p.hexSignature() == p2.hexSignature() {
		t.Errorf("Got the same signature with two different secret keys")
	}

	err = ValidatePayload(p.payload, headerWithRolledKey, p.secret)
	if err != nil {
		t.Errorf("Expected to be able to decode webhook with old key after rolling key, but got %v", err)
	}
	err = ValidatePayload(p.payload, headerWithRolledKey, p2.secret)
	if err != nil {
		t.Errorf("Expected to be able to decode webhook with new key after rolling key, but got %v", err)
	}

	p = newSignedPayload(func(p *signedPayload) {
		p.timestamp = time.Now().Add(-15 * time.Second)
	})
	err = ValidatePayloadWithTolerance(p.payload, p.header, p.secret, 10*time.Second)
	if err != ErrTooOld {
		t.Errorf("Received %v error when validating timestamp outside of allowed timing window", err)
	}

	err = ValidatePayloadWithTolerance(p.payload, p.header, p.secret, 20*time.Second)
	if err != nil {
		t.Errorf("Received %v error when validating timestamp inside allowed timing window", err)
	}

	p = newSignedPayload(func(p *signedPayload) {
		p.timestamp = time.Unix(12345, 0)
	})
	err = ValidatePayloadIgnoringTolerance(p.payload, p.header, p.secret)
	if err != nil {
		t.Errorf("Received %v error when timestamp outside window but no tolerance specified", err)
	}
}
